<script setup lang="ts">
import type { ButtonConfig } from '@/types/flow-preview'
import { ExternalLink } from 'lucide-vue-next'

defineProps<{
  buttons: ButtonConfig[]
  disabled?: boolean
}>()

const emit = defineEmits<{
  select: [button: ButtonConfig]
}>()

function handleClick(button: ButtonConfig) {
  if (button.type === 'url' && button.url) {
    // For URL buttons, just acknowledge in preview
    emit('select', button)
  } else {
    emit('select', button)
  }
}
</script>

<template>
  <div class="mt-1 space-y-1">
    <button
      v-for="btn in buttons"
      :key="btn.id"
      class="w-full bg-white dark:bg-[#1a1a1a] text-black dark:text-white text-sm font-medium py-2.5 px-4 rounded-lg shadow-sm border border-zinc-200 dark:border-zinc-800 flex items-center justify-center gap-1.5 transition-colors"
      :class="{
        'hover:bg-zinc-50 dark:hover:bg-zinc-900 cursor-pointer': !disabled,
        'opacity-50 cursor-not-allowed': disabled
      }"
      :disabled="disabled"
      @click="handleClick(btn)"
    >
      <ExternalLink v-if="btn.type === 'url'" class="h-4 w-4" />
      {{ btn.title || 'Option' }}
    </button>
  </div>
</template>
