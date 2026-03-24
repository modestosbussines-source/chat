<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Search,
  Filter,
  MoreVertical,
  RefreshCw,
  Archive,
  Trash2,
  Star,
  Phone,
  Video,
  Users,
  Settings
} from 'lucide-vue-next'
import Button from '@/components/ui/button/Button.vue'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import type { Conversation } from '@/types/inbox'

interface Props {
  conversation?: Conversation | null
  showSearch?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  conversation: null,
  showSearch: true
})

const emit = defineEmits<{
  (e: 'search', query: string): void
  (e: 'refresh'): void
  (e: 'archive'): void
  (e: 'delete'): void
  (e: 'star'): void
  (e: 'call'): void
  (e: 'video-call'): void
  (e: 'assign'): void
  (e: 'settings'): void
}>()

const { t } = useI18n()

const contact = computed(() => props.conversation?.contact)
const isAssigned = computed(() => !!props.conversation?.assignedTo)
const isStarred = computed(() => props.conversation?.isStarred)

const contactInitials = computed(() => {
  const name = contact.value?.name || ''
  if (!name) return '?'
  return name
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
})

const statusLabel = computed(() => {
  switch (props.conversation?.status) {
    case 'open':
      return t('inbox.status.open')
    case 'pending':
      return t('inbox.status.pending')
    case 'resolved':
      return t('inbox.status.resolved')
    default:
      return ''
  }
})

const statusColor = computed(() => {
  switch (props.conversation?.status) {
    case 'open':
      return 'bg-success'
    case 'pending':
      return 'bg-warning'
    case 'resolved':
      return 'bg-muted-foreground'
    default:
      return 'bg-muted-foreground'
  }
})
</script>

<template>
  <div class="inbox-header flex items-center justify-between px-4 py-3 border-b border-border bg-card">
    <!-- Left: Contact Info -->
    <div v-if="conversation && contact" class="flex items-center gap-3">
      <!-- Avatar -->
      <div class="relative">
        <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-sm font-medium text-primary">
          {{ contactInitials }}
        </div>
        <div
          v-if="contact.isOnline"
          class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-success border-2 border-card"
        />
      </div>
      
      <!-- Contact Details -->
      <div>
        <h2 class="text-sm font-semibold text-foreground flex items-center gap-2">
          {{ contact.name || 'Unknown' }}
          <span
            v-if="isStarred"
            class="text-warning"
          >
            <Star class="w-3.5 h-3.5 fill-current" />
          </span>
        </h2>
        <div class="flex items-center gap-2 text-xs text-muted-foreground">
          <span class="flex items-center gap-1">
            <span :class="['w-2 h-2 rounded-full', statusColor]" />
            {{ statusLabel }}
          </span>
          <span v-if="isAssigned">·</span>
          <span v-if="isAssigned" class="text-primary">
            {{ conversation?.assignedTo?.name }}
          </span>
        </div>
      </div>
    </div>
    
    <!-- Center: No conversation selected -->
    <div v-else class="flex items-center gap-3">
      <div class="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
        <Users class="w-5 h-5 text-muted-foreground" />
      </div>
      <div>
        <h2 class="text-sm font-semibold text-foreground">{{ t('inbox.header.title') }}</h2>
        <p class="text-xs text-muted-foreground">{{ t('inbox.header.selectConversation') }}</p>
      </div>
    </div>
    
    <!-- Right: Actions -->
    <div class="flex items-center gap-2">
      <!-- Action Buttons (when conversation selected) -->
      <template v-if="conversation">
        <!-- Call -->
        <Button
          v-if="contact?.phone"
          variant="ghost"
          size="icon"
          class="h-8 w-8"
          @click="emit('call')"
          :title="t('inbox.actions.call')"
        >
          <Phone class="w-4 h-4" />
        </Button>
        
        <!-- Video Call -->
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8"
          @click="emit('video-call')"
          :title="t('inbox.actions.videoCall')"
        >
          <Video class="w-4 h-4" />
        </Button>
        
        <!-- Assign -->
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8"
          @click="emit('assign')"
          :title="t('inbox.actions.assign')"
        >
          <Users class="w-4 h-4" />
        </Button>
      </template>
      
      <!-- More Actions Dropdown -->
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" size="icon" class="h-8 w-8">
            <MoreVertical class="w-4 h-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-48">
          <DropdownMenuItem @click="emit('refresh')">
            <RefreshCw class="w-4 h-4 mr-2" />
            {{ t('inbox.actions.refresh') }}
          </DropdownMenuItem>
          
          <template v-if="conversation">
            <DropdownMenuItem @click="emit('star')">
              <Star class="w-4 h-4 mr-2" />
              {{ isStarred ? t('inbox.actions.unstar') : t('inbox.actions.star') }}
            </DropdownMenuItem>
            
            <DropdownMenuItem @click="emit('archive')">
              <Archive class="w-4 h-4 mr-2" />
              {{ t('inbox.actions.archive') }}
            </DropdownMenuItem>
            
            <DropdownMenuSeparator />
            
            <DropdownMenuItem @click="emit('delete')" class="text-destructive focus:text-destructive">
              <Trash2 class="w-4 h-4 mr-2" />
              {{ t('inbox.actions.delete') }}
            </DropdownMenuItem>
          </template>
          
          <DropdownMenuSeparator />
          
          <DropdownMenuItem @click="emit('settings')">
            <Settings class="w-4 h-4 mr-2" />
            {{ t('inbox.actions.settings') }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </div>
</template>

<style scoped>
.inbox-header {
  min-height: 60px;
}
</style>
