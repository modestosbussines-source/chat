<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatDistanceToNow } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { 
  Wifi, 
  WifiOff, 
  Loader2, 
  MoreVertical, 
  Trash2, 
  Power, 
  QrCode,
  Phone,
  User,
  Clock
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import Avatar from '@/components/ui/avatar/Avatar.vue'
import AvatarFallback from '@/components/ui/avatar/AvatarFallback.vue'
import AvatarImage from '@/components/ui/avatar/AvatarImage.vue'
import type { EvolutionInstance } from '@/types/evolution'

interface Props {
  instance: EvolutionInstance
  isActive?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isActive: false
})

const emit = defineEmits<{
  (e: 'connect', instance: EvolutionInstance): void
  (e: 'disconnect', instance: EvolutionInstance): void
  (e: 'delete', instance: EvolutionInstance): void
  (e: 'view-details', instance: EvolutionInstance): void
}>()

const { t } = useI18n()

const statusConfig = computed(() => {
  switch (props.instance.status) {
    case 'connected':
      return {
        label: t('evolution.status.connected'),
        color: 'bg-success',
        textColor: 'text-success',
        icon: Wifi,
        animate: ''
      }
    case 'connecting':
      return {
        label: t('evolution.status.connecting'),
        color: 'bg-warning',
        textColor: 'text-warning',
        icon: Loader2,
        animate: 'animate-spin'
      }
    case 'disconnected':
    default:
      return {
        label: t('evolution.status.disconnected'),
        color: 'bg-muted-foreground',
        textColor: 'text-muted-foreground',
        icon: WifiOff,
        animate: ''
      }
  }
})

const lastActivity = computed(() => {
  if (!props.instance.last_activity) return null
  return formatDistanceToNow(new Date(props.instance.last_activity), {
    addSuffix: true,
    locale: ptBR
  })
})

const connectedTime = computed(() => {
  if (!props.instance.connected_at) return null
  return formatDistanceToNow(new Date(props.instance.connected_at), {
    addSuffix: false,
    locale: ptBR
  })
})

const instanceInitials = computed(() => {
  const name = props.instance.display_name || props.instance.instance_name
  return name
    .slice(0, 2)
    .toUpperCase()
})

const handleConnect = () => {
  emit('connect', props.instance)
}

const handleDisconnect = () => {
  emit('disconnect', props.instance)
}

const handleDelete = () => {
  emit('delete', props.instance)
}

const handleViewDetails = () => {
  emit('view-details', props.instance)
}
</script>

<template>
  <div
    :class="[
      'instance-card group relative rounded-xl border border-border bg-card overflow-hidden',
      'transition-all duration-300 ease-out',
      'hover:shadow-lg hover:shadow-primary/5 hover:border-primary/30',
      'hover:-translate-y-1',
      isActive && 'ring-2 ring-primary/30 border-primary'
    ]"
  >
    <!-- Status Bar -->
    <div
      :class="[
        'absolute top-0 left-0 right-0 h-1 transition-colors duration-300',
        statusConfig.color
      ]"
    />

    <div class="p-4 pt-5">
      <!-- Header -->
      <div class="flex items-start justify-between gap-3 mb-4">
        <div class="flex items-center gap-3 min-w-0">
          <!-- Avatar -->
          <div class="relative flex-shrink-0">
            <Avatar class="w-12 h-12 ring-2 ring-border transition-all duration-300 group-hover:ring-primary/30">
              <AvatarImage :src="instance.profile_pic || undefined" :alt="instance.display_name" />
              <AvatarFallback class="bg-gradient-to-br from-primary/20 to-secondary/20 text-primary font-semibold">
                {{ instanceInitials }}
              </AvatarFallback>
            </Avatar>
            
            <!-- Status Badge -->
            <div
              :class="[
                'absolute -bottom-1 -right-1 w-5 h-5 rounded-full border-2 border-card flex items-center justify-center',
                statusConfig.color
              ]"
            >
              <component
                :is="statusConfig.icon"
                :class="['w-3 h-3 text-white', statusConfig.animate]"
              />
            </div>
          </div>

          <!-- Info -->
          <div class="min-w-0 flex-1">
            <h3 class="font-semibold text-foreground truncate text-sm">
              {{ instance.display_name || instance.instance_name }}
            </h3>
            <p class="text-xs text-muted-foreground truncate">
              {{ instance.instance_name }}
            </p>
          </div>
        </div>

        <!-- Actions -->
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity"
            >
              <MoreVertical class="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-48">
            <DropdownMenuItem @click="handleViewDetails">
              <User class="w-4 h-4 mr-2" />
              {{ t('evolution.actions.details') }}
            </DropdownMenuItem>
            
            <template v-if="instance.status === 'connected'">
              <DropdownMenuItem @click="handleDisconnect" class="text-warning focus:text-warning">
                <Power class="w-4 h-4 mr-2" />
                {{ t('evolution.actions.disconnect') }}
              </DropdownMenuItem>
            </template>
            <template v-else>
              <DropdownMenuItem @click="handleConnect">
                <QrCode class="w-4 h-4 mr-2" />
                {{ t('evolution.actions.connect') }}
              </DropdownMenuItem>
            </template>

            <DropdownMenuSeparator />

            <DropdownMenuItem @click="handleDelete" class="text-destructive focus:text-destructive">
              <Trash2 class="w-4 h-4 mr-2" />
              {{ t('evolution.actions.delete') }}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <!-- Status Info -->
      <div class="space-y-2">
        <!-- Status Badge -->
        <div class="flex items-center gap-2">
          <span
            :class="[
              'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium',
              instance.status === 'connected' 
                ? 'bg-success/10 text-success'
                : instance.status === 'connecting'
                  ? 'bg-warning/10 text-warning'
                  : 'bg-muted text-muted-foreground'
            ]"
          >
            <span :class="['w-1.5 h-1.5 rounded-full', statusConfig.color]" />
            {{ statusConfig.label }}
          </span>
          
          <span v-if="connectedTime" class="text-xs text-muted-foreground">
            {{ connectedTime }}
          </span>
        </div>

        <!-- Phone & Profile (when connected) -->
        <div v-if="instance.status === 'connected'" class="space-y-1.5">
          <div v-if="instance.phone" class="flex items-center gap-2 text-xs text-muted-foreground">
            <Phone class="w-3 h-3" />
            <span>{{ instance.phone }}</span>
          </div>
          <div v-if="instance.profile_name" class="flex items-center gap-2 text-xs text-muted-foreground">
            <User class="w-3 h-3" />
            <span>{{ instance.profile_name }}</span>
          </div>
        </div>

        <!-- Last Activity -->
        <div v-if="lastActivity" class="flex items-center gap-2 text-xs text-muted-foreground">
          <Clock class="w-3 h-3" />
          <span>{{ t('evolution.lastActivity') }}: {{ lastActivity }}</span>
        </div>

        <!-- Error Message -->
        <div
          v-if="instance.last_error"
          class="mt-2 p-2 rounded-lg bg-destructive/10 text-destructive text-xs"
        >
          {{ instance.last_error }}
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mt-4 flex gap-2">
        <Button
          v-if="instance.status !== 'connected'"
          size="sm"
          class="flex-1 gap-2"
          @click="handleConnect"
        >
          <QrCode class="w-4 h-4" />
          {{ t('evolution.actions.connect') }}
        </Button>
        
        <Button
          v-else
          size="sm"
          variant="outline"
          class="flex-1 gap-2"
          @click="handleDisconnect"
        >
          <Power class="w-4 h-4" />
          {{ t('evolution.actions.disconnect') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.instance-card {
  position: relative;
}

/* Pulse animation for connecting status */
.instance-card:has([class*="animate-spin"])::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent, hsl(var(--warning) / 0.05), transparent);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0; }
  50% { opacity: 1; }
}

/* Touch device optimizations */
@media (hover: none) {
  .instance-card:hover {
    transform: none;
  }
  
  .instance-card .group-hover\:opacity-100 {
    opacity: 1;
  }
  
  .instance-card:active {
    transform: scale(0.98);
  }
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .instance-card {
    transition: none;
  }
  
  .instance-card:hover {
    transform: none;
  }
}
</style>
