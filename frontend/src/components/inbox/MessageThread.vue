<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { format, isToday, isYesterday } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import {
  Check,
  CheckCheck,
  Clock,
  AlertCircle,
  Image,
  File,
  Mic,
  Video,
  MapPin,
  Reply,
  Smile,
  MoreHorizontal
} from 'lucide-vue-next'
import Avatar from '@/components/ui/avatar/Avatar.vue'
import AvatarFallback from '@/components/ui/avatar/AvatarFallback.vue'
import AvatarImage from '@/components/ui/avatar/AvatarImage.vue'
import Button from '@/components/ui/button/Button.vue'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import type { Message } from '@/types/inbox'

interface Props {
  messages: Message[]
  contactName?: string
  contactAvatar?: string
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  contactName: 'Contact',
  loading: false
})

const emit = defineEmits<{
  (e: 'reply', message: Message): void
  (e: 'react', message: Message, emoji: string): void
  (e: 'delete', message: Message): void
  (e: 'retry', message: Message): void
}>()

const { t } = useI18n()

const messagesContainer = ref<HTMLElement | null>(null)
const hoveredMessageId = ref<string | null>(null)
const newlyAddedMessageId = ref<string | null>(null)

// Auto scroll to bottom on new messages
watch(() => props.messages.length, (newLength, oldLength) => {
  if (newLength > oldLength) {
    // Mark new message for animation
    const newMessage = props.messages[props.messages.length - 1]
    if (newMessage) {
      newlyAddedMessageId.value = newMessage.id
      setTimeout(() => {
        newlyAddedMessageId.value = null
      }, 500)
    }
    
    nextTick(() => {
      scrollToBottom()
    })
  }
})

onMounted(() => {
  scrollToBottom()
})

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTo({
      top: messagesContainer.value.scrollHeight,
      behavior: 'smooth'
    })
  }
}

// Group messages by date
const groupedMessages = computed(() => {
  const groups: { date: string; messages: Message[] }[] = []
  let currentDate = ''

  props.messages.forEach((message) => {
    const messageDate = formatMessageDate(message.createdAt)
    
    if (messageDate !== currentDate) {
      currentDate = messageDate
      groups.push({ date: messageDate, messages: [message] })
    } else {
      groups[groups.length - 1].messages.push(message)
    }
  })

  return groups
})

const formatMessageDate = (dateString: string): string => {
  const date = new Date(dateString)
  
  if (isToday(date)) {
    return t('inbox.chat.today')
  }
  
  if (isYesterday(date)) {
    return t('inbox.chat.yesterday')
  }
  
  return format(date, "d 'de' MMMM", { locale: ptBR })
}

const formatMessageTime = (dateString: string): string => {
  return format(new Date(dateString), 'HH:mm')
}

const isOwnMessage = (message: Message): boolean => {
  return message.senderType === 'agent'
}

const getMessageStatus = (message: Message): string => {
  if (!isOwnMessage(message)) return ''
  
  switch (message.status) {
    case 'sent':
      return 'sent'
    case 'delivered':
      return 'delivered'
    case 'read':
      return 'read'
    case 'failed':
      return 'failed'
    default:
      return 'sending'
  }
}

const getMediaIcon = (type: string) => {
  switch (type) {
    case 'image':
      return Image
    case 'audio':
      return Mic
    case 'video':
      return Video
    case 'location':
      return MapPin
    default:
      return File
  }
}

const handleReply = (message: Message) => {
  emit('reply', message)
}

const handleReact = (message: Message, emoji: string) => {
  emit('react', message, emoji)
}

const handleDelete = (message: Message) => {
  emit('delete', message)
}

const handleRetry = (message: Message) => {
  emit('retry', message)
}

// Quick reactions
const quickReactions = ['👍', '❤️', '😂', '😮', '😢', '🙏']
</script>

<template>
  <div class="message-thread flex flex-col h-full bg-chat-background">
    <!-- Messages Container -->
    <div
      ref="messagesContainer"
      class="flex-1 overflow-y-auto px-3 md:px-4 py-4 smooth-scroll"
    >
      <!-- Loading State -->
      <div v-if="loading" class="space-y-4">
        <div v-for="i in 5" :key="i" :class="['flex', i % 2 === 0 ? 'justify-end' : 'justify-start']">
          <div class="skeleton rounded-2xl animate-pulse" :style="{ width: i % 2 === 0 ? '60%' : '75%', height: '48px' }" />
        </div>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="messages.length === 0"
        class="flex flex-col items-center justify-center h-full text-center animate-fade-in"
      >
        <div class="w-16 h-16 mb-4 rounded-full bg-muted flex items-center justify-center animate-float">
          <svg
            class="w-8 h-8 text-muted-foreground"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
          </svg>
        </div>
        <p class="text-sm font-medium text-foreground animate-fade-in-up stagger-1">
          {{ t('inbox.chat.empty.title') }}
        </p>
        <p class="text-xs text-muted-foreground mt-1 animate-fade-in-up stagger-2">
          {{ t('inbox.chat.empty.description') }}
        </p>
      </div>

      <!-- Messages -->
      <div v-else class="space-y-4">
        <div 
          v-for="(group, groupIndex) in groupedMessages" 
          :key="group.date"
          :class="['animate-fade-in', `stagger-${Math.min(groupIndex + 1, 8)}`]"
        >
          <!-- Date Separator -->
          <div class="flex items-center justify-center my-4">
            <div class="px-3 py-1 rounded-full bg-muted text-xs text-muted-foreground font-medium">
              {{ group.date }}
            </div>
          </div>

          <!-- Messages in this group -->
          <TransitionGroup name="message" tag="div" class="space-y-2">
            <div
              v-for="(message, index) in group.messages"
              :key="message.id"
              :class="[
                'flex gap-2 group relative',
                isOwnMessage(message) ? 'flex-row-reverse' : '',
                newlyAddedMessageId === message.id && 'animate-bubble-appear'
              ]"
              :style="{ animationDelay: `${index * 50}ms` }"
              @mouseenter="hoveredMessageId = message.id"
              @mouseleave="hoveredMessageId = null"
            >
              <!-- Avatar (for incoming) -->
              <Avatar
                v-if="!isOwnMessage"
                class="w-8 h-8 flex-shrink-0 mt-1 transition-transform duration-200 hover:scale-110"
              >
                <AvatarImage :src="contactAvatar" :alt="contactName" />
                <AvatarFallback class="text-xs bg-primary/10 text-primary">
                  {{ contactName?.charAt(0)?.toUpperCase() || '?' }}
                </AvatarFallback>
              </Avatar>

              <!-- Message Bubble -->
              <div :class="['max-w-[85%] sm:max-w-[70%] flex flex-col', isOwnMessage(message) ? 'items-end' : 'items-start']">
                <!-- Reply Preview -->
                <div
                  v-if="message.replyTo"
                  class="reply-preview text-xs max-w-full mb-1 animate-fade-in-down"
                >
                  <p class="font-medium">{{ message.replyTo.senderName }}</p>
                  <p class="truncate">{{ message.replyTo.content }}</p>
                </div>

                <!-- Media Content -->
                <div
                  v-if="message.media"
                  :class="[
                    'chat-bubble transition-all duration-200',
                    isOwnMessage(message) ? 'chat-bubble-outgoing' : 'chat-bubble-incoming'
                  ]"
                >
                  <div class="flex items-center gap-2">
                    <component :is="getMediaIcon(message.media.type)" class="w-4 h-4" />
                    <span class="text-sm">{{ message.media.name || 'Media' }}</span>
                  </div>
                </div>

                <!-- Text Content -->
                <div
                  v-if="message.content"
                  :class="[
                    'chat-bubble relative transition-all duration-200',
                    isOwnMessage(message) ? 'chat-bubble-outgoing' : 'chat-bubble-incoming',
                    newlyAddedMessageId === message.id && 'message-bubble'
                  ]"
                >
                  <p class="whitespace-pre-wrap break-words">{{ message.content }}</p>
                  
                  <!-- Time & Status -->
                  <div class="chat-bubble-time">
                    <span>{{ formatMessageTime(message.createdAt) }}</span>
                    
                    <!-- Status Icon (for own messages) -->
                    <template v-if="isOwnMessage(message)">
                      <Check
                        v-if="getMessageStatus(message) === 'sent'"
                        class="status-icon w-3 h-3"
                      />
                      <CheckCheck
                        v-else-if="getMessageStatus(message) === 'delivered'"
                        class="status-icon w-3 h-3"
                      />
                      <CheckCheck
                        v-else-if="getMessageStatus(message) === 'read'"
                        class="status-icon w-3 h-3 text-blue-400"
                      />
                      <AlertCircle
                        v-else-if="getMessageStatus(message) === 'failed'"
                        class="status-icon w-3 h-3 text-destructive"
                      />
                      <Clock
                        v-else
                        class="status-icon w-3 h-3 animate-pulse"
                      />
                    </template>
                  </div>
                </div>

                <!-- Reactions -->
                <div
                  v-if="message.reactions && message.reactions.length > 0"
                  class="reactions-display flex gap-1 mt-1 animate-scale-in"
                >
                  <span
                    v-for="reaction in message.reactions"
                    :key="reaction.id"
                    class="reaction-badge text-xs transition-transform hover:scale-110"
                  >
                    {{ reaction.emoji }} {{ reaction.count > 1 ? reaction.count : '' }}
                  </span>
                </div>

                <!-- Quick Actions (on hover - desktop only) -->
                <div
                  v-if="hoveredMessageId === message.id"
                  :class="[
                    'absolute top-0 hidden sm:flex items-center gap-1 animate-fade-in',
                    isOwnMessage(message) ? '-left-24' : '-right-24'
                  ]"
                >
                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button variant="ghost" size="icon" class="h-7 w-7 bg-card shadow-sm rounded-full">
                        <Smile class="w-3.5 h-3.5" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                      <div class="flex gap-1 p-1">
                        <button
                          v-for="emoji in quickReactions"
                          :key="emoji"
                          class="p-1.5 hover:bg-muted rounded transition-all duration-150 hover:scale-125"
                          @click="handleReact(message, emoji)"
                        >
                          {{ emoji }}
                        </button>
                      </div>
                    </DropdownMenuContent>
                  </DropdownMenu>

                  <Button
                    variant="ghost"
                    size="icon"
                    class="h-7 w-7 bg-card shadow-sm rounded-full"
                    @click="handleReply(message)"
                  >
                    <Reply class="w-3.5 h-3.5" />
                  </Button>

                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button variant="ghost" size="icon" class="h-7 w-7 bg-card shadow-sm rounded-full">
                        <MoreHorizontal class="w-3.5 h-3.5" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                      <DropdownMenuItem @click="handleDelete(message)">
                        {{ t('inbox.chat.actions.delete') }}
                      </DropdownMenuItem>
                      <DropdownMenuItem v-if="message.status === 'failed'" @click="handleRetry(message)">
                        {{ t('inbox.chat.actions.retry') }}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
            </div>
          </TransitionGroup>
        </div>
      </div>
    </div>

    <!-- Typing Indicator (optional) -->
    <!-- <div class="px-4 py-2">
      <div class="typing-indicator">
        <span></span>
        <span></span>
        <span></span>
      </div>
    </div> -->
  </div>
</template>

<style scoped>
.message-thread {
  background-color: hsl(var(--chat-bg));
}

/* Message animations */
.message-enter-active {
  animation: messageAppear 0.35s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.message-leave-active {
  animation: messageDisappear 0.2s ease-in forwards;
}

.message-move {
  transition: transform 0.3s ease;
}

@keyframes messageAppear {
  from {
    opacity: 0;
    transform: scale(0.92) translateY(8px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

@keyframes messageDisappear {
  to {
    opacity: 0;
    transform: scale(0.92) translateY(-8px);
  }
}

/* Bubble appear animation */
.message-bubble {
  animation: bubblePop 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes bubblePop {
  0% {
    opacity: 0;
    transform: scale(0.8);
  }
  50% {
    transform: scale(1.02);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

/* Scale animation */
.animate-scale-in {
  animation: scaleIn 0.3s ease-out forwards;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* Smooth scroll */
.smooth-scroll {
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
}

/* Typing indicator */
.typing-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 12px 16px;
  background: hsl(var(--muted));
  border-radius: 1rem;
  width: fit-content;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: hsl(var(--muted-foreground));
  border-radius: 50%;
  animation: typingDot 1.4s infinite ease-in-out both;
}

.typing-indicator span:nth-child(1) { animation-delay: 0s; }
.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes typingDot {
  0%, 80%, 100% {
    transform: translateY(0);
    opacity: 0.4;
  }
  40% {
    transform: translateY(-6px);
    opacity: 1;
  }
}

/* Mobile optimizations */
@media (max-width: 640px) {
  .chat-bubble {
    max-width: 88% !important;
  }
}

/* Reduce motion for accessibility */
@media (prefers-reduced-motion: reduce) {
  .message-enter-active,
  .message-leave-active,
  .message-bubble,
  .animate-scale-in {
    animation: none !important;
  }
  
  .smooth-scroll {
    scroll-behavior: auto;
  }
}
</style>
