<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useContactsStore } from '@/stores/contacts'
import { usersService, chatbotService } from '@/services/api'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { Badge } from '@/components/ui/badge'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/components/ui/popover'
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from '@/components/ui/alert-dialog'
import { LogOut, User, Settings, Moon, Sun, Globe } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { getInitials, getAvatarGradient } from '@/lib/utils'
import { useColorMode } from '@/composables/useColorMode'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'

const { t } = useI18n()
const { isDark, toggleTheme } = useColorMode()

defineProps<{
  collapsed?: boolean
}>()

const emit = defineEmits<{
  logout: []
}>()

const authStore = useAuthStore()
const contactsStore = useContactsStore()
const isUserMenuOpen = ref(false)
const isUpdatingAvailability = ref(false)
const isCheckingTransfers = ref(false)
const showAwayWarning = ref(false)
const awayWarningTransferCount = ref(0)

const handleAvailabilityChange = async (checked: boolean) => {
  if (!checked) {
    isCheckingTransfers.value = true
    try {
      const response = await chatbotService.listTransfers({ status: 'active' })
      const data = response.data.data || response.data
      const transfers = data.transfers || []
      const userId = authStore.user?.id
      const myActiveTransfers = transfers.filter((t: any) => t.agent_id === userId)

      if (myActiveTransfers.length > 0) {
        awayWarningTransferCount.value = myActiveTransfers.length
        showAwayWarning.value = true
        return
      }
    } catch (error) {
      console.error('Failed to check transfers:', error)
    } finally {
      isCheckingTransfers.value = false
    }
  }

  await setAvailability(checked)
}

const confirmGoAway = async () => {
  showAwayWarning.value = false
  await setAvailability(false)
}

const setAvailability = async (checked: boolean) => {
  isUpdatingAvailability.value = true
  try {
    const response = await usersService.updateAvailability(checked)
    const data = response.data.data
    authStore.setAvailability(checked, data.break_started_at)

    if (checked) {
      toast.success(t('userMenu.available'), {
        description: t('userMenu.availableDesc')
      })
    } else {
      const transfersReturned = data.transfers_to_queue || 0
      toast.success(t('userMenu.away'), {
        description: transfersReturned > 0
          ? t('userMenu.transfersReturned', { count: transfersReturned })
          : t('userMenu.awayDesc')
      })

      if (transfersReturned > 0) {
        contactsStore.fetchContacts()
      }
    }
  } catch (error) {
    toast.error(t('common.error'), {
      description: t('userMenu.failedUpdateAvailability')
    })
  } finally {
    isUpdatingAvailability.value = false
  }
}

const breakDuration = ref('')
let breakTimerInterval: ReturnType<typeof setInterval> | null = null

const updateBreakDuration = () => {
  if (!authStore.breakStartedAt) {
    breakDuration.value = ''
    return
  }
  const start = new Date(authStore.breakStartedAt)
  const now = new Date()
  const diffMs = now.getTime() - start.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const hours = Math.floor(diffMins / 60)
  const mins = diffMins % 60

  if (hours > 0) {
    breakDuration.value = `${hours}h ${mins}m`
  } else {
    breakDuration.value = `${mins}m`
  }
}

watch(() => authStore.isAvailable, (available) => {
  if (!available && authStore.breakStartedAt) {
    updateBreakDuration()
    breakTimerInterval = setInterval(updateBreakDuration, 60000)
  } else if (breakTimerInterval) {
    clearInterval(breakTimerInterval)
    breakTimerInterval = null
    breakDuration.value = ''
  }
}, { immediate: true })

onMounted(() => {
  authStore.restoreBreakTime()
  if (!authStore.isAvailable && authStore.breakStartedAt) {
    updateBreakDuration()
    breakTimerInterval = setInterval(updateBreakDuration, 60000)
  }
})

onUnmounted(() => {
  if (breakTimerInterval) {
    clearInterval(breakTimerInterval)
  }
})

const handleLogout = () => {
  emit('logout')
}
</script>

<template>
  <div class="border-t border-border p-3">
    <Popover v-model:open="isUserMenuOpen">
      <PopoverTrigger as-child>
        <Button
          variant="ghost"
          :class="[
            'flex items-center justify-start w-full h-auto px-2 py-2 gap-3 hover:bg-muted rounded-lg',
            collapsed && 'md:justify-center'
          ]"
          aria-label="User menu"
        >
          <div class="relative">
            <Avatar class="h-9 w-9 ring-2 ring-border">
              <AvatarImage :src="undefined" />
              <AvatarFallback :class="'text-sm bg-gradient-to-br text-white ' + getAvatarGradient(authStore.user?.full_name || 'U')">
                {{ getInitials(authStore.user?.full_name || 'U') }}
              </AvatarFallback>
            </Avatar>
            <div
              :class="[
                'absolute -bottom-0.5 -right-0.5 w-3 h-3 rounded-full border-2 border-card',
                authStore.isAvailable ? 'bg-success' : 'bg-muted-foreground'
              ]"
            />
          </div>
          <div v-if="!collapsed" class="flex flex-col items-start text-left flex-1 min-w-0">
            <span class="text-sm font-medium truncate max-w-[140px] text-foreground">
              {{ authStore.user?.full_name }}
            </span>
            <span class="text-xs text-muted-foreground truncate max-w-[140px]">
              {{ authStore.user?.email }}
            </span>
          </div>
        </Button>
      </PopoverTrigger>
      <PopoverContent side="top" align="start" class="w-64 p-2">
        <!-- User Info -->
        <div class="flex items-center gap-3 px-2 py-3 border-b border-border mb-2">
          <Avatar class="h-10 w-10">
            <AvatarImage :src="undefined" />
            <AvatarFallback :class="'text-sm bg-gradient-to-br text-white ' + getAvatarGradient(authStore.user?.full_name || 'U')">
              {{ getInitials(authStore.user?.full_name || 'U') }}
            </AvatarFallback>
          </Avatar>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-foreground truncate">
              {{ authStore.user?.full_name }}
            </p>
            <p class="text-xs text-muted-foreground truncate">
              {{ authStore.user?.email }}
            </p>
          </div>
        </div>

        <!-- Availability Toggle -->
        <div class="px-2 py-2 mb-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-sm text-foreground">{{ $t('userMenu.status') }}</span>
              <Badge
                :class="authStore.isAvailable 
                  ? 'bg-success/10 text-success' 
                  : 'bg-muted text-muted-foreground'"
              >
                {{ authStore.isAvailable ? $t('userMenu.available') : $t('userMenu.away') }}
              </Badge>
            </div>
            <Switch
              :checked="authStore.isAvailable"
              :disabled="isUpdatingAvailability || isCheckingTransfers"
              aria-label="Toggle availability status"
              @update:checked="handleAvailabilityChange"
            />
          </div>
          <p v-if="!authStore.isAvailable && breakDuration" class="text-xs text-muted-foreground mt-1">
            Ausente há {{ breakDuration }}
          </p>
        </div>

        <Separator class="my-2" />

        <!-- Menu Items -->
        <div class="space-y-1">
          <RouterLink to="/profile" @click="isUserMenuOpen = false">
            <Button
              variant="ghost"
              class="w-full justify-start px-2 py-2 h-auto text-sm font-normal"
            >
              <User class="mr-2 h-4 w-4 text-muted-foreground" />
              <span>{{ $t('userMenu.profile') }}</span>
            </Button>
          </RouterLink>

          <Button
            variant="ghost"
            class="w-full justify-start px-2 py-2 h-auto text-sm font-normal"
            @click="toggleTheme"
          >
            <Sun v-if="isDark" class="mr-2 h-4 w-4 text-muted-foreground" />
            <Moon v-else class="mr-2 h-4 w-4 text-muted-foreground" />
            <span>{{ isDark ? 'Modo claro' : 'Modo escuro' }}</span>
          </Button>

          <div class="px-2 py-1">
            <div class="flex items-center gap-2 text-sm text-muted-foreground">
              <Globe class="h-4 w-4" />
              <span>{{ $t('userMenu.language') }}</span>
            </div>
            <div class="mt-1">
              <LanguageSwitcher />
            </div>
          </div>
        </div>

        <Separator class="my-2" />

        <!-- Logout -->
        <Button
          variant="ghost"
          class="w-full justify-start px-2 py-2 h-auto text-sm font-normal text-destructive hover:text-destructive hover:bg-destructive/10"
          @click="handleLogout"
        >
          <LogOut class="mr-2 h-4 w-4" />
          <span>{{ $t('userMenu.logOut') }}</span>
        </Button>
      </PopoverContent>
    </Popover>
  </div>

  <!-- Away Warning Dialog -->
  <AlertDialog :open="showAwayWarning">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ $t('userMenu.awayWarningTitle') }}</AlertDialogTitle>
        <AlertDialogDescription>
          {{ $t('userMenu.awayWarningDesc', { count: awayWarningTransferCount }) }}
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <Button variant="outline" @click="showAwayWarning = false">{{ $t('common.cancel') }}</Button>
        <Button @click="confirmGoAway" :disabled="isUpdatingAvailability">{{ $t('userMenu.goAway') }}</Button>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
