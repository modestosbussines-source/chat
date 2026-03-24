<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Copy, Check, AlertCircle, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { toast } from 'vue-sonner'

interface Props {
  qrCode?: string | null
  loading?: boolean
  error?: string | null
  autoRefresh?: boolean
  refreshInterval?: number
}

const props = withDefaults(defineProps<Props>(), {
  qrCode: null,
  loading: false,
  error: null,
  autoRefresh: true,
  refreshInterval: 15000 // 15 seconds
})

const emit = defineEmits<{
  (e: 'refresh'): void
}>()

const { t } = useI18n()

const isRefreshing = ref(false)
const timeRemaining = ref(props.refreshInterval / 1000)
const copied = ref(false)
let refreshInterval: ReturnType<typeof setInterval> | null = null
let countdownInterval: ReturnType<typeof setInterval> | null = null

// Progress percentage for the circular timer
const progressPercentage = computed(() => {
  return (timeRemaining.value / (props.refreshInterval / 1000)) * 100
})

// Watch for QR code changes to reset timer
watch(() => props.qrCode, () => {
  timeRemaining.value = props.refreshInterval / 1000
})

// Auto-refresh logic
const startAutoRefresh = () => {
  if (!props.autoRefresh) return
  
  // Clear existing intervals
  stopAutoRefresh()
  
  // Countdown timer
  countdownInterval = setInterval(() => {
    timeRemaining.value--
    if (timeRemaining.value <= 0) {
      timeRemaining.value = props.refreshInterval / 1000
    }
  }, 1000)
  
  // Refresh trigger
  refreshInterval = setInterval(() => {
    handleRefresh()
  }, props.refreshInterval)
}

const stopAutoRefresh = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
  if (countdownInterval) {
    clearInterval(countdownInterval)
    countdownInterval = null
  }
}

const handleRefresh = async () => {
  isRefreshing.value = true
  timeRemaining.value = props.refreshInterval / 1000
  
  try {
    emit('refresh')
  } finally {
    setTimeout(() => {
      isRefreshing.value = false
    }, 500)
  }
}

const copyQRCode = async () => {
  if (!props.qrCode) return
  
  try {
    // Copy the base64 data URL
    await navigator.clipboard.writeText(props.qrCode)
    copied.value = true
    toast.success(t('evolution.qrCode.copied'))
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    toast.error(t('evolution.qrCode.copyFailed'))
  }
}

// Lifecycle
onMounted(() => {
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<template>
  <div class="qr-code-display flex flex-col items-center">
    <!-- QR Code Container -->
    <div class="relative w-64 h-64 rounded-2xl border-2 border-dashed border-border bg-muted/30 flex items-center justify-center overflow-hidden">
      <!-- Loading State -->
      <div v-if="loading || isRefreshing" class="absolute inset-0 flex items-center justify-center bg-muted/50 backdrop-blur-sm animate-fade-in">
        <Loader2 class="w-10 h-10 text-primary animate-spin" />
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="absolute inset-0 flex flex-col items-center justify-center p-6 text-center animate-fade-in">
        <div class="w-12 h-12 rounded-full bg-destructive/10 flex items-center justify-center mb-3">
          <AlertCircle class="w-6 h-6 text-destructive" />
        </div>
        <p class="text-sm text-destructive font-medium mb-1">{{ t('evolution.qrCode.error') }}</p>
        <p class="text-xs text-muted-foreground">{{ error }}</p>
        <Button
          variant="outline"
          size="sm"
          class="mt-4 gap-2"
          @click="handleRefresh"
        >
          <RefreshCw class="w-4 h-4" />
          {{ t('evolution.qrCode.retry') }}
        </Button>
      </div>

      <!-- QR Code Image -->
      <div v-else-if="qrCode" class="animate-scale-in">
        <img
          :src="qrCode"
          :alt="t('evolution.qrCode.alt')"
          class="w-full h-full object-contain"
        />
      </div>

      <!-- Empty State -->
      <div v-else class="flex flex-col items-center justify-center text-center p-6">
        <div class="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-3">
          <svg class="w-8 h-8 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
          </svg>
        </div>
        <p class="text-sm text-muted-foreground">
          {{ t('evolution.qrCode.empty') }}
        </p>
      </div>

      <!-- Circular Timer -->
      <div
        v-if="qrCode && !loading && !error"
        class="absolute top-3 right-3 w-10 h-10"
      >
        <svg class="w-full h-full -rotate-90" viewBox="0 0 36 36">
          <!-- Background circle -->
          <circle
            cx="18"
            cy="18"
            r="15"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            class="text-border"
          />
          <!-- Progress circle -->
          <circle
            cx="18"
            cy="18"
            r="15"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-dasharray="94.25"
            :stroke-dashoffset="94.25 - (progressPercentage / 100) * 94.25"
            stroke-linecap="round"
            class="text-primary transition-all duration-1000"
          />
        </svg>
        <span class="absolute inset-0 flex items-center justify-center text-[10px] font-medium text-muted-foreground">
          {{ timeRemaining }}s
        </span>
      </div>

      <!-- Refresh Button (overlay) -->
      <button
        v-if="qrCode && !loading"
        @click="handleRefresh"
        :disabled="isRefreshing"
        class="absolute bottom-3 right-3 p-2 rounded-lg bg-card/80 backdrop-blur-sm border border-border shadow-sm hover:bg-card transition-all hover:scale-105 active:scale-95 disabled:opacity-50"
      >
        <RefreshCw :class="['w-4 h-4 text-muted-foreground', isRefreshing && 'animate-spin']" />
      </button>
    </div>

    <!-- Instructions -->
    <div class="mt-6 text-center max-w-xs">
      <h4 class="text-sm font-semibold text-foreground mb-2">
        {{ t('evolution.qrCode.instructionTitle') }}
      </h4>
      <ol class="text-xs text-muted-foreground space-y-1.5">
        <li class="flex items-center gap-2">
          <span class="w-5 h-5 rounded-full bg-primary/10 text-primary text-[10px] font-bold flex items-center justify-center">1</span>
          {{ t('evolution.qrCode.step1') }}
        </li>
        <li class="flex items-center gap-2">
          <span class="w-5 h-5 rounded-full bg-primary/10 text-primary text-[10px] font-bold flex items-center justify-center">2</span>
          {{ t('evolution.qrCode.step2') }}
        </li>
        <li class="flex items-center gap-2">
          <span class="w-5 h-5 rounded-full bg-primary/10 text-primary text-[10px] font-bold flex items-center justify-center">3</span>
          {{ t('evolution.qrCode.step3') }}
        </li>
      </ol>
    </div>

    <!-- Actions -->
    <div class="mt-4 flex gap-2">
      <Button
        variant="outline"
        size="sm"
        class="gap-2"
        @click="copyQRCode"
        :disabled="!qrCode"
      >
        <component :is="copied ? Check : Copy" class="w-4 h-4" />
        {{ copied ? t('evolution.qrCode.copied') : t('evolution.qrCode.copy') }}
      </Button>
      
      <Button
        variant="outline"
        size="sm"
        class="gap-2"
        @click="handleRefresh"
        :disabled="isRefreshing"
      >
        <RefreshCw :class="['w-4 h-4', isRefreshing && 'animate-spin']" />
        {{ t('evolution.qrCode.refresh') }}
      </Button>
    </div>

    <!-- Auto-refresh indicator -->
    <p v-if="qrCode && !error" class="mt-3 text-xs text-muted-foreground">
      {{ t('evolution.qrCode.autoRefresh', { seconds: timeRemaining }) }}
    </p>
  </div>
</template>

<style scoped>
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

/* Fade animation */
.animate-fade-in {
  animation: fadeIn 0.2s ease-out forwards;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Touch device optimizations */
@media (hover: none) {
  button:active {
    transform: scale(0.95);
  }
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .animate-scale-in,
  .animate-fade-in {
    animation: none !important;
  }
  
  * {
    transition-duration: 0.01ms !important;
  }
}
</style>
