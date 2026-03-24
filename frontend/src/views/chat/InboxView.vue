<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useContactsStore, type Contact, type Message } from '@/stores/contacts'
import { useAuthStore } from '@/stores/auth'
import { useTagsStore } from '@/stores/tags'
import { wsService } from '@/services/websocket'
import { contactsService, messagesService } from '@/services/api'
import { toast } from 'vue-sonner'

// OMNI Design System Components
import {
  ConversationList,
  ConversationItem,
  ContactCard,
  InboxHeader,
  MessageThread,
  MessageInput
} from '@/components/inbox'
import SkeletonLoader from '@/components/ui/skeleton/SkeletonLoader.vue'

// Existing UI components
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useColorMode } from '@/composables/useColorMode'
import { ArrowLeft, PanelLeftClose, PanelLeft } from 'lucide-vue-next'

const { t } = useI18n()
const contactsStore = useContactsStore()
const authStore = useAuthStore()
const tagsStore = useTagsStore()
const { isDark, toggleTheme } = useColorMode()

// State
const activeConversationId = ref<string | null>(null)
const activeFilter = ref<'all' | 'unread' | 'assigned' | 'unassigned'>('all')
const isLoading = ref(false)
const isAssignDialogOpen = ref(false)
const isInfoPanelOpen = ref(false) // Closed by default on mobile
const replyToMessage = ref<any>(null)
const isMobileConversationsOpen = ref(true) // Show conversations list on mobile
const isDragging = ref(false)

// Show info panel only on desktop
const showInfoPanel = computed(() => isInfoPanelOpen.value && activeConversationId.value)

// Transform contacts to conversations format
const conversations = computed(() => {
  return contactsStore.contacts.map(contact => ({
    id: contact.id,
    contact: {
      id: contact.id,
      name: contact.name || contact.phone_number,
      phone: contact.phone_number,
      avatar: contact.avatar_url,
      isOnline: false,
      company: null,
      email: null,
      location: null,
      tags: contact.tags?.map(tag => ({ id: tag, name: tag })) || [],
      createdAt: contact.created_at,
      notes: null
    },
    lastMessage: contact.last_message_at ? {
      id: `${contact.id}-last`,
      content: contact.last_message_content || '',
      createdAt: contact.last_message_at,
      senderType: 'customer'
    } : null,
    unreadCount: contact.unread_count || 0,
    status: 'open',
    channel: 'whatsapp',
    tags: contact.tags?.map(tag => ({ id: tag, name: tag })) || [],
    assignedTo: contact.assigned_to ? { id: contact.assigned_to, name: contact.assigned_to } : null,
    isStarred: false
  }))
})

// Active conversation
const activeConversation = computed(() => {
  if (!activeConversationId.value) return null
  return conversations.value.find(c => c.id === activeConversationId.value)
})

// Messages transformed
const transformedMessages = computed(() => {
  return contactsStore.messages.map(msg => ({
    id: msg.id,
    content: msg.content,
    createdAt: msg.created_at,
    senderType: msg.direction === 'outgoing' ? 'agent' : 'customer',
    senderName: msg.direction === 'outgoing' ? 'Você' : (contactsStore.currentContact?.name || 'Cliente'),
    status: msg.status || 'sent',
    media: msg.media_url ? {
      type: msg.message_type === 'image' ? 'image' : msg.message_type === 'audio' ? 'audio' : msg.message_type === 'video' ? 'video' : 'file',
      url: msg.media_url,
      name: msg.media_filename || 'Media'
    } : null,
    replyTo: msg.is_reply && msg.reply_to_message ? {
      id: msg.reply_to_message_id,
      content: msg.reply_to_message.content,
      senderName: msg.reply_to_message.direction === 'outgoing' ? 'Você' : (contactsStore.currentContact?.name || 'Cliente')
    } : null,
    reactions: []
  }))
})

// Handlers
const handleSelectConversation = async (conversation: any) => {
  activeConversationId.value = conversation.id
  const contact = contactsStore.contacts.find(c => c.id === conversation.id)
  if (contact) {
    await contactsStore.setCurrentContact(contact)
  }
  // On mobile, show chat view
  isMobileConversationsOpen.value = false
}

const handleFilterChange = (filter: string) => {
  activeFilter.value = filter as any
}

const handleSendMessage = async (content: string, attachments: File[]) => {
  if (!contactsStore.currentContact || !content.trim()) return
  
  try {
    await contactsStore.sendMessage(content, attachments)
    replyToMessage.value = null
  } catch (error) {
    toast.error(t('chat.errors.sendMessageFailed'))
  }
}

const handleReply = (message: any) => {
  replyToMessage.value = message
}

const handleReact = (message: any, emoji: string) => {
  toast.success(`${emoji} ${t('chat.reactions.added')}`)
}

const handleDeleteMessage = async (message: any) => {
  toast.success(t('chat.messages.deleted'))
}

const handleRetryMessage = async (message: any) => {
  // Retry logic
}

const handleCancelReply = () => {
  replyToMessage.value = null
}

const handleCall = () => {
  if (contactsStore.currentContact?.phone_number) {
    window.open(`tel:${contactsStore.currentContact.phone_number}`)
  }
}

const handleVideoCall = () => {
  toast.info(t('chat.actions.videoCallComingSoon'))
}

const handleAssign = () => {
  isAssignDialogOpen.value = true
}

const handleStar = () => {
  toast.success(t('chat.actions.conversationStarred'))
}

const handleArchive = () => {
  toast.success(t('chat.actions.conversationArchived'))
}

const handleDeleteConversation = () => {
  toast.success(t('chat.actions.conversationDeleted'))
}

const handleSettings = () => {
  // Navigate to settings
}

const handleRefresh = async () => {
  isLoading.value = true
  try {
    await contactsStore.fetchContacts()
    toast.success(t('chat.actions.refreshed'))
  } catch (error) {
    toast.error(t('chat.errors.refreshFailed'))
  } finally {
    isLoading.value = false
  }
}

const handleBackToList = () => {
  isMobileConversationsOpen.value = true
  activeConversationId.value = null
}

const toggleInfoPanel = () => {
  isInfoPanelOpen.value = !isInfoPanelOpen.value
}

// Lifecycle
onMounted(async () => {
  isLoading.value = true
  try {
    await contactsStore.fetchContacts()
  } catch (error) {
    toast.error(t('chat.errors.loadContactsFailed'))
  } finally {
    isLoading.value = false
  }
})

// Watch for contact changes
watch(() => contactsStore.currentContact, (newContact) => {
  if (newContact) {
    activeConversationId.value = newContact.id
  }
})
</script>

<template>
  <div class="inbox-view flex h-full bg-background">
    <!-- Mobile Layout -->
    <div class="flex-1 flex md:hidden">
      <!-- Conversations List (Mobile) -->
      <Transition name="slide-right" mode="out-in">
        <div
          v-if="isMobileConversationsOpen"
          key="conversations"
          class="flex-1 flex flex-col"
        >
          <InboxHeader
            :conversation="null"
            @refresh="handleRefresh"
            @settings="handleSettings"
          />
          
          <div class="flex-1 overflow-hidden">
            <ConversationList
              :conversations="conversations"
              :active-conversation-id="activeConversationId || undefined"
              :loading="isLoading"
              :filter="activeFilter"
              @select="handleSelectConversation"
              @filter-change="handleFilterChange"
            />
          </div>
        </div>

        <!-- Chat View (Mobile) -->
        <div
          v-else
          key="chat"
          class="flex-1 flex flex-col"
        >
          <!-- Chat Header with Back Button -->
          <div class="flex items-center gap-2 px-3 py-2 border-b border-border bg-card">
            <Button
              variant="ghost"
              size="icon"
              class="h-9 w-9 flex-shrink-0"
              @click="handleBackToList"
            >
              <ArrowLeft class="h-5 w-5" />
            </Button>
            
            <div v-if="activeConversation?.contact" class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-medium text-primary">
                  {{ activeConversation.contact.name?.charAt(0)?.toUpperCase() || '?' }}
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-foreground truncate">
                    {{ activeConversation.contact.name }}
                  </p>
                  <p class="text-xs text-muted-foreground">
                    <span class="inline-block w-2 h-2 rounded-full bg-success mr-1" />
                    {{ $t('chat.online') }}
                  </p>
                </div>
              </div>
            </div>
            
            <Button
              variant="ghost"
              size="icon"
              class="h-9 w-9"
              @click="toggleInfoPanel"
            >
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </Button>
          </div>

          <!-- Messages -->
          <div class="flex-1 min-h-0">
            <MessageThread
              :messages="transformedMessages"
              :contact-name="activeConversation?.contact?.name"
              :contact-avatar="activeConversation?.contact?.avatar"
              :loading="isLoading"
              @reply="handleReply"
              @react="handleReact"
              @delete="handleDeleteMessage"
              @retry="handleRetryMessage"
            />
          </div>

          <!-- Message Input -->
          <MessageInput
            :placeholder="t('inbox.input.placeholder')"
            :reply-to="replyToMessage"
            @send="handleSendMessage"
            @cancel-reply="handleCancelReply"
          />
        </div>
      </Transition>
    </div>

    <!-- Desktop Layout -->
    <div class="hidden md:flex flex-1">
      <!-- Conversations Sidebar -->
      <div 
        :class="[
          'flex flex-col border-r border-border bg-card transition-all duration-300',
          showInfoPanel ? 'w-80 xl:w-96' : 'w-80 xl:w-96 flex-1 max-w-md'
        ]"
      >
        <InboxHeader
          :conversation="null"
          @refresh="handleRefresh"
          @settings="handleSettings"
        />
        
        <div class="flex-1 overflow-hidden">
          <ConversationList
            :conversations="conversations"
            :active-conversation-id="activeConversationId || undefined"
            :loading="isLoading"
            :filter="activeFilter"
            @select="handleSelectConversation"
            @filter-change="handleFilterChange"
          />
        </div>
      </div>

      <!-- Chat Area -->
      <div class="flex-1 flex flex-col min-w-0">
        <!-- No Conversation Selected -->
        <div
          v-if="!activeConversationId"
          class="flex-1 flex items-center justify-center animate-fade-in"
        >
          <div class="text-center">
            <div class="w-20 h-20 rounded-2xl bg-primary/10 flex items-center justify-center mx-auto mb-4 animate-float">
              <svg
                class="w-10 h-10 text-primary"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-foreground mb-1 animate-fade-in-up stagger-1">
              {{ t('inbox.empty.title') }}
            </h3>
            <p class="text-sm text-muted-foreground animate-fade-in-up stagger-2">
              {{ t('inbox.empty.description') }}
            </p>
          </div>
        </div>

        <!-- Chat Interface -->
        <template v-else>
          <InboxHeader
            :conversation="activeConversation"
            @call="handleCall"
            @video-call="handleVideoCall"
            @assign="handleAssign"
            @star="handleStar"
            @archive="handleArchive"
            @delete="handleDeleteConversation"
            @settings="handleSettings"
            @refresh="handleRefresh"
          />

          <div class="flex-1 min-h-0">
            <MessageThread
              :messages="transformedMessages"
              :contact-name="activeConversation?.contact?.name"
              :contact-avatar="activeConversation?.contact?.avatar"
              :loading="isLoading"
              @reply="handleReply"
              @react="handleReact"
              @delete="handleDeleteMessage"
              @retry="handleRetryMessage"
            />
          </div>

          <MessageInput
            :placeholder="t('inbox.input.placeholder')"
            :reply-to="replyToMessage"
            @send="handleSendMessage"
            @cancel-reply="handleCancelReply"
          />
        </template>
      </div>

      <!-- Contact Info Panel (Desktop) -->
      <Transition name="slide-left">
        <div
          v-if="showInfoPanel"
          class="w-80 xl:w-96 border-l border-border bg-card overflow-y-auto"
        >
          <ContactCard
            :contact="activeConversation!.contact!"
            @message="() => {}"
            @call="handleCall"
            @edit="() => {}"
          />
        </div>
      </Transition>
    </div>

    <!-- Mobile Info Panel -->
    <Transition name="slide-up">
      <div
        v-if="isInfoPanelOpen && activeConversation && !isMobileConversationsOpen"
        class="fixed inset-0 z-50 md:hidden bg-background"
      >
        <div class="flex flex-col h-full">
          <!-- Header -->
          <div class="flex items-center gap-3 px-4 py-3 border-b border-border">
            <Button
              variant="ghost"
              size="icon"
              class="h-9 w-9"
              @click="isInfoPanelOpen = false"
            >
              <ArrowLeft class="h-5 w-5" />
            </Button>
            <h2 class="text-lg font-semibold">{{ t('chat.contactInfo') }}</h2>
          </div>
          
          <!-- Content -->
          <div class="flex-1 overflow-y-auto">
            <ContactCard
              :contact="activeConversation.contact!"
              @message="() => { isInfoPanelOpen = false }"
              @call="handleCall"
              @edit="() => {}"
            />
          </div>
        </div>
      </div>
    </Transition>

    <!-- Assign Dialog -->
    <Dialog v-model:open="isAssignDialogOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ t('chat.assignDialog.title') }}</DialogTitle>
          <DialogDescription>
            {{ t('chat.assignDialog.description') }}
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  </div>
</template>

<style scoped>
/* Mobile conversation/chat transitions */
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s ease;
}

.slide-right-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-right-leave-to {
  transform: translateX(-30%);
  opacity: 0;
}

/* Info panel slide */
.slide-left-enter-active,
.slide-left-leave-active {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s ease;
}

.slide-left-enter-from,
.slide-left-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

/* Mobile slide up for info panel */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s ease;
}

.slide-up-enter-from {
  transform: translateY(100%);
  opacity: 0;
}

.slide-up-leave-to {
  transform: translateY(100%);
  opacity: 0;
}

/* Animation utilities */
.animate-fade-in {
  animation: fadeIn 0.3s ease-out forwards;
}

.animate-fade-in-up {
  animation: fadeInUp 0.4s ease-out forwards;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

/* Stagger delays */
.stagger-1 { animation-delay: 0.1s; }
.stagger-2 { animation-delay: 0.2s; }
.stagger-3 { animation-delay: 0.3s; }
.stagger-4 { animation-delay: 0.4s; }
</style>
