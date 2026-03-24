<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Send,
  Paperclip,
  Smile,
  Image,
  Mic,
  X,
  Reply,
  AtSign,
  Hash,
  Zap,
  MessageSquare,
  Video,
  FileIcon
} from 'lucide-vue-next'
import Button from '@/components/ui/button/Button.vue'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import type { Message } from '@/types/inbox'

interface Props {
  placeholder?: string
  disabled?: boolean
  replyTo?: Message | null
  showEmojiPicker?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Mensagem...',
  disabled: false,
  replyTo: null,
  showEmojiPicker: true
})

const emit = defineEmits<{
  (e: 'send', content: string, attachments: File[]): void
  (e: 'typing'): void
  (e: 'cancel-reply'): void
  (e: 'attach', files: File[]): void
  (e: 'open-omnis'): void
}>()

const { t } = useI18n()

const message = ref('')
const attachments = ref<File[]>([])
const fileInput = ref<HTMLInputElement | null>(null)
const textareaRef = ref<HTMLTextAreaElement | null>(null)

const canSend = computed(() => {
  return message.value.trim().length > 0 || attachments.value.length > 0
})

const characterCount = computed(() => message.value.length)
const maxCharacters = 4096 // WhatsApp limit

const isOverLimit = computed(() => characterCount.value > maxCharacters)

// Auto-resize textarea
watch(message, () => {
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto'
    textareaRef.value.style.height = Math.min(textareaRef.value.scrollHeight, 120) + 'px'
  }
})

const handleInput = () => {
  emit('typing')
}

const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    handleSend()
  }
}

const handleSend = () => {
  if (!canSend.value || isOverLimit.value) return
  
  emit('send', message.value.trim(), [...attachments.value])
  
  // Reset
  message.value = ''
  attachments.value = []
}

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files) {
    const files = Array.from(input.files)
    attachments.value = [...attachments.value, ...files]
    emit('attach', files)
  }
  // Reset input
  input.value = ''
}

const openFilePicker = () => {
  fileInput.value?.click()
}

const removeAttachment = (index: number) => {
  attachments.value.splice(index, 1)
}

const clearReply = () => {
  emit('cancel-reply')
}

const insertEmoji = (emoji: string) => {
  message.value += emoji
}

// Common emojis for quick access
const quickEmojis = ['😊', '😂', '❤️', '👍', '🙏', '🎉', '😍', '🤣', '😢', '😎']
</script>

<template>
  <div class="message-input-container border-t border-border bg-card">
    <!-- Reply Preview -->
    <div
      v-if="replyTo"
      class="flex items-center gap-2 px-4 py-2 bg-muted border-l-2 border-primary"
    >
      <Reply class="w-4 h-4 text-primary flex-shrink-0" />
      <div class="flex-1 min-w-0">
        <p class="text-xs font-medium text-primary">{{ replyTo.senderName }}</p>
        <p class="text-xs text-muted-foreground truncate">{{ replyTo.content }}</p>
      </div>
      <Button
        variant="ghost"
        size="icon"
        class="h-6 w-6 flex-shrink-0"
        @click="clearReply"
      >
        <X class="w-4 h-4" />
      </Button>
    </div>

    <!-- Attachments Preview -->
    <div
      v-if="attachments.length > 0"
      class="flex gap-2 px-4 py-2 overflow-x-auto"
    >
      <div
        v-for="(file, index) in attachments"
        :key="index"
        class="relative flex items-center gap-2 px-3 py-2 bg-muted rounded-lg"
      >
        <Paperclip class="w-4 h-4 text-muted-foreground flex-shrink-0" />
        <span class="text-sm truncate max-w-[150px]">{{ file.name }}</span>
        <Button
          variant="ghost"
          size="icon"
          class="h-5 w-5 flex-shrink-0 -mr-1"
          @click="removeAttachment(index)"
        >
          <X class="w-3 h-3" />
        </Button>
      </div>
    </div>

    <!-- Input Area -->
    <div class="flex items-end gap-2 px-4 py-3">
      <!-- Attachment & Omnis Buttons -->
      <div class="flex gap-1">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="ghost" size="icon" class="h-9 w-9 flex-shrink-0">
              <Paperclip class="w-5 h-5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start">
            <DropdownMenuItem @click="openFilePicker">
              <File class="w-4 h-4 mr-2" />
              {{ t('inbox.input.attachments.file') }}
            </DropdownMenuItem>
            <DropdownMenuItem @click="openFilePicker">
              <Image class="w-4 h-4 mr-2" />
              {{ t('inbox.input.attachments.image') }}
            </DropdownMenuItem>
            <DropdownMenuItem @click="openFilePicker">
              <Mic class="w-4 h-4 mr-2" />
              {{ t('inbox.input.attachments.audio') }}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <!-- Omnis Button -->
        <Button
          variant="ghost"
          size="icon"
          class="h-9 w-9 flex-shrink-0"
          @click="emit('open-omnis')"
          title="Mensagens Rápidas"
        >
          <Zap class="w-5 h-5" />
        </Button>
      </div>
        </DropdownMenu>
      </div>

      <!-- Hidden File Input -->
      <input
        ref="fileInput"
        type="file"
        multiple
        class="hidden"
        @change="handleFileSelect"
      />

      <!-- Text Input -->
      <div class="flex-1 relative">
        <textarea
          ref="textareaRef"
          v-model="message"
          :placeholder="placeholder"
          :disabled="disabled"
          rows="1"
          class="w-full px-4 py-2.5 text-sm bg-muted border border-transparent rounded-xl resize-none focus:outline-none focus:border-ring focus:bg-background transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          @input="handleInput"
          @keydown="handleKeyDown"
        />
        
        <!-- Character Counter -->
        <div
          v-if="characterCount > 0"
          :class="[
            'absolute -bottom-4 right-2 text-[10px]',
            isOverLimit ? 'text-destructive' : 'text-muted-foreground'
          ]"
        >
          {{ characterCount }}/{{ maxCharacters }}
        </div>
      </div>

      <!-- Emoji Button -->
      <DropdownMenu v-if="showEmojiPicker">
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" size="icon" class="h-9 w-9 flex-shrink-0">
            <Smile class="w-5 h-5" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-64">
          <div class="grid grid-cols-5 gap-1 p-2">
            <button
              v-for="emoji in quickEmojis"
              :key="emoji"
              class="p-2 text-xl hover:bg-muted rounded transition-colors"
              @click="insertEmoji(emoji)"
            >
              {{ emoji }}
            </button>
          </div>
        </DropdownMenuContent>
      </DropdownMenu>

      <!-- Send Button -->
      <Button
        :disabled="!canSend || isOverLimit"
        size="icon"
        class="h-9 w-9 flex-shrink-0"
        @click="handleSend"
      >
        <Send class="w-4 h-4" />
      </Button>
    </div>
  </div>
</template>

<style scoped>
.message-input-container {
  position: relative;
}
</style>
