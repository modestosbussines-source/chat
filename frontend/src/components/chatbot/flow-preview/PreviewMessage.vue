<script setup lang="ts">
import { computed } from 'vue'
import type { SimulationMessage } from '@/types/flow-preview'
import { Bug, Info } from 'lucide-vue-next'

const props = defineProps<{
  message: SimulationMessage
}>()

const formattedTime = computed(() => {
  return props.message.timestamp.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  })
})

const isBot = computed(() => props.message.type === 'bot')
const isUser = computed(() => props.message.type === 'user')
const isSystem = computed(() => props.message.type === 'system')
const isDebug = computed(() => props.message.type === 'debug')
</script>

<template>
  <!-- Bot Message -->
  <div v-if="isBot" class="flex justify-start">
    <div class="max-w-[85%]">
      <div
        class="bg-white dark:bg-black rounded-lg rounded-tl-none shadow-sm p-3"
        :class="{ 'border-l-2 border-red-400': message.isValidationError, 'border border-zinc-200 dark:border-zinc-800': !message.isValidationError }"
      >
        <p class="text-sm text-black dark:text-white whitespace-pre-wrap">
          {{ message.content }}
        </p>
        <p class="text-[10px] text-zinc-400 text-right mt-1">{{ formattedTime }}</p>
      </div>

      <!-- Show step name for debugging -->
      <p v-if="message.stepName" class="text-[10px] text-gray-400 mt-0.5 ml-1">
        Step: {{ message.stepName }}
      </p>
    </div>
  </div>

  <!-- User Message -->
  <div v-else-if="isUser" class="flex justify-end">
    <div class="max-w-[85%]">
      <div class="bg-black dark:bg-white rounded-lg rounded-tr-none shadow-sm p-3 border border-zinc-800 dark:border-zinc-200">
        <p class="text-sm text-white dark:text-black">{{ message.content }}</p>
        <p class="text-[10px] text-zinc-400 text-right mt-1 flex items-center justify-end gap-1">
          {{ formattedTime }}
          <svg class="h-4 w-4 text-zinc-400" viewBox="0 0 24 24" fill="currentColor">
            <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z"/>
          </svg>
        </p>
      </div>
    </div>
  </div>

  <!-- System Message -->
  <div v-else-if="isSystem" class="flex justify-center">
    <div class="bg-zinc-100 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-xs text-zinc-700 dark:text-zinc-400 px-3 py-1.5 rounded-lg flex items-center gap-1.5">
      <Info class="h-3 w-3" />
      <span>{{ message.content }}</span>
    </div>
  </div>

  <!-- Debug Message -->
  <div v-else-if="isDebug" class="flex justify-center">
    <div class="bg-zinc-100 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-xs text-zinc-700 dark:text-zinc-400 px-3 py-1.5 rounded-lg flex items-center gap-1.5 max-w-[90%]">
      <Bug class="h-3 w-3 flex-shrink-0" />
      <span class="break-all">{{ message.content }}</span>
    </div>
  </div>
</template>
