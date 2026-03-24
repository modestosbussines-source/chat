<script setup lang="ts">
import { computed } from 'vue'
import { formatDistanceToNow } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import Avatar from '@/components/ui/avatar/Avatar.vue'
import AvatarFallback from '@/components/ui/avatar/AvatarFallback.vue'
import AvatarImage from '@/components/ui/avatar/AvatarImage.vue'
import type { Conversation } from '@/types/inbox'

interface Props {
  conversation: Conversation
  active?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  active: false
})

const emit = defineEmits<{
  (e: 'click'): void
}>()

const contact = computed(() => props.conversation.contact)
const lastMessage = computed(() => props.conversation.lastMessage)
const unreadCount = computed(() => props.conversation.unreadCount)

const formattedTime = computed(() => {
  if (!lastMessage.value?.createdAt) return ''
  return formatDistanceToNow(new Date(lastMessage.value.createdAt), {
    addSuffix: false,
    locale: ptBR
  })
})

const truncatedMessage = computed(() => {
  if (!lastMessage.value?.content) return ''
  const content = lastMessage.value.content
  return content.length > 45 ? content.slice(0, 45) + '...' : content
})

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

const statusColor = computed(() => {
  switch (props.conversation.status) {
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

const channelIcon = computed(() => {
  switch (props.conversation.channel) {
    case 'whatsapp':
      return '💬'
    case 'instagram':
      return '📷'
    case 'facebook':
      return '📘'
    case 'telegram':
      return '✈️'
    case 'email':
      return '📧'
    default:
      return '💬'
  }
})
</script>

<template>
  <div
    @click="emit('click')"
    :class="[
      'conversation-item mx-2 rounded-xl cursor-pointer',
      'transition-all duration-200 ease-out',
      'hover:bg-muted/80 active:scale-[0.98] active:bg-muted',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/30 focus-visible:ring-offset-1',
      active && 'bg-primary/10 border-l-2 border-primary',
      unreadCount > 0 && 'font-medium'
    ]"
    role="button"
    tabindex="0"
    @keydown.enter="emit('click')"
    @keydown.space.prevent="emit('click')"
  >
    <div class="flex items-center gap-3 p-3">
      <!-- Avatar -->
      <div class="relative flex-shrink-0">
        <Avatar class="w-11 h-11 transition-transform duration-200 group-hover:scale-105">
          <AvatarImage :src="contact?.avatar" :alt="contact?.name" />
          <AvatarFallback 
            :class="[
              'text-sm font-medium',
              active ? 'bg-primary text-primary-foreground' : 'bg-primary/10 text-primary'
            ]"
          >
            {{ contactInitials }}
          </AvatarFallback>
        </Avatar>
        
        <!-- Online Status Indicator with pulse -->
        <div
          v-if="contact?.isOnline"
          class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-success border-2 border-card animate-pulse-ring"
        />
        
        <!-- Channel Icon -->
        <div class="absolute -bottom-0.5 -right-0.5 text-xs bg-card rounded-full p-0.5">
          {{ channelIcon }}
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center justify-between gap-2 mb-0.5">
          <span 
            :class="[
              'text-sm truncate transition-colors duration-150',
              active ? 'text-primary font-semibold' : 'text-foreground'
            ]"
          >
            {{ contact?.name || 'Unknown' }}
          </span>
          <span class="flex-shrink-0 text-[11px] text-muted-foreground">
            {{ formattedTime }}
          </span>
        </div>
        
        <div class="flex items-center justify-between gap-2">
          <p 
            :class="[
              'text-xs truncate transition-colors duration-150',
              unreadCount > 0 ? 'text-foreground' : 'text-muted-foreground'
            ]"
          >
            <span
              v-if="lastMessage?.senderType === 'agent'"
              :class="['font-medium mr-1', active ? 'text-primary' : 'text-primary/80']"
            >
              Você:
            </span>
            {{ truncatedMessage || 'Nenhuma mensagem' }}
          </p>
          
          <!-- Unread Badge with animation -->
          <div
            v-if="unreadCount > 0"
            class="unread-badge flex-shrink-0 badge-pulse"
          >
            {{ unreadCount > 99 ? '99+' : unreadCount }}
          </div>
          
          <!-- Status Indicator -->
          <div
            v-else
            :class="['w-2 h-2 rounded-full flex-shrink-0 transition-colors duration-200', statusColor]"
          />
        </div>
        
        <!-- Tags Preview with animation -->
        <div
          v-if="conversation.tags && conversation.tags.length > 0"
          class="flex gap-1 mt-1.5 overflow-hidden"
        >
          <span
            v-for="(tag, index) in conversation.tags.slice(0, 2)"
            :key="tag.id"
            :class="[
              'px-1.5 py-0.5 text-[10px] rounded-md bg-muted text-muted-foreground',
              'transition-all duration-200 hover:bg-muted-foreground/10',
              'animate-fade-in'
            ]"
            :style="{ animationDelay: `${index * 50}ms` }"
          >
            {{ tag.name }}
          </span>
          <span
            v-if="conversation.tags.length > 2"
            class="px-1.5 py-0.5 text-[10px] rounded-md bg-muted text-muted-foreground"
          >
            +{{ conversation.tags.length - 2 }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.conversation-item {
  position: relative;
  overflow: hidden;
}

/* Unread badge pulse animation */
.badge-pulse {
  animation: badgePulse 2s ease-in-out infinite;
}

@keyframes badgePulse {
  0%, 100% {
    box-shadow: 0 0 0 0 hsl(var(--primary) / 0.3);
  }
  50% {
    box-shadow: 0 0 0 6px hsl(var(--primary) / 0);
  }
}

/* Active state indicator slide */
.conversation-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: hsl(var(--primary));
  border-radius: 0 2px 2px 0;
  animation: slideIn 0.2s ease-out;
}

@keyframes slideIn {
  from {
    transform: scaleY(0);
    opacity: 0;
  }
  to {
    transform: scaleY(1);
    opacity: 1;
  }
}

/* Hover highlight effect */
.conversation-item::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent, hsl(var(--primary) / 0.03), transparent);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.conversation-item:hover::after {
  opacity: 1;
}

/* Animation utilities */
.animate-fade-in {
  animation: fadeIn 0.2s ease-out forwards;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Pulse ring for online indicator */
.animate-pulse-ring {
  animation: pulseRing 2s ease-in-out infinite;
}

@keyframes pulseRing {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.8;
  }
}

/* Touch device optimizations */
@media (hover: none) {
  .conversation-item:active {
    transform: scale(0.98);
    background: hsl(var(--muted));
  }
  
  .conversation-item:active .unread-badge {
    animation: none;
  }
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .conversation-item,
  .badge-pulse,
  .animate-pulse-ring {
    animation: none !important;
    transition: none !important;
  }
  
  .conversation-item::after {
    display: none;
  }
}
</style>
