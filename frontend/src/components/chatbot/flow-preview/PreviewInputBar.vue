<script setup lang="ts">
import { ref, computed } from 'vue'
import { Send, Mic } from 'lucide-vue-next'

const props = defineProps<{
  inputType: string | null
  disabled?: boolean
  placeholder?: string
}>()

const emit = defineEmits<{
  submit: [value: string]
}>()

const inputValue = ref('')

const inputPlaceholder = computed(() => {
  if (props.placeholder) return props.placeholder

  switch (props.inputType) {
    case 'email':
      return 'Enter your email address...'
    case 'phone':
      return 'Enter your phone number...'
    case 'number':
      return 'Enter a number...'
    case 'date':
      return 'Enter a date (YYYY-MM-DD)...'
    case 'text':
      return 'Type a message...'
    default:
      return 'Type a message...'
  }
})

const inputTypeAttr = computed(() => {
  switch (props.inputType) {
    case 'email':
      return 'email'
    case 'phone':
      return 'tel'
    case 'number':
      return 'number'
    case 'date':
      return 'date'
    default:
      return 'text'
  }
})

const isEnabled = computed(() => {
  return !props.disabled && props.inputType && props.inputType !== 'none' && props.inputType !== 'button'
})

// Convert input to string (handles number inputs)
const inputString = computed(() => String(inputValue.value ?? '').trim())

function handleSubmit() {
  if (!inputString.value || !isEnabled.value) return

  emit('submit', inputString.value)
  inputValue.value = ''
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    handleSubmit()
  }
}
</script>

<template>
  <div class="bg-zinc-100 dark:bg-zinc-900 px-3 py-2 flex items-center gap-2 border-t border-zinc-200 dark:border-zinc-800">
    <div class="flex-1 bg-white dark:bg-zinc-800 rounded-full px-4 py-2 border border-zinc-200 dark:border-zinc-700">
      <input
        v-if="isEnabled"
        v-model="inputValue"
        :type="inputTypeAttr"
        :placeholder="inputPlaceholder"
        class="w-full text-sm bg-transparent border-none outline-none text-black dark:text-white placeholder:text-zinc-400"
        @keydown="handleKeydown"
      />
      <p v-else class="text-sm text-zinc-400">
        {{ disabled ? 'Waiting...' : 'Type a message' }}
      </p>
    </div>

    <button
      class="w-10 h-10 rounded-full flex items-center justify-center transition-colors"
      :class="{
        'bg-black dark:bg-white cursor-pointer': isEnabled && inputString,
        'bg-zinc-300 dark:bg-zinc-600 cursor-not-allowed': !isEnabled || !inputString
      }"
      :disabled="!isEnabled || !inputString"
      @click="handleSubmit"
    >
      <Send v-if="inputString" class="h-5 w-5 text-white dark:text-black" />
      <Mic v-else class="h-5 w-5 text-zinc-500 dark:text-zinc-400" />
    </button>
  </div>
</template>
