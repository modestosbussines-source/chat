<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { X, Smartphone, Wifi, WifiOff, Loader2, Check, AlertTriangle } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import QRCodeDisplay from './QRCodeDisplay.vue'
import type { EvolutionInstance, EvolutionInstanceStatus } from '@/types/evolution'

interface Props {
  instance: EvolutionInstance | null
  open: boolean
  loading?: boolean
  qrCode?: string | null
  error?: string | null
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  qrCode: null,
  error: null
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'refresh-qr'): void
  (e: 'disconnect'): void
}>()

const { t } = useI18n()

const isPolling = ref(false)
let pollingInterval: ReturnType<typeof setInterval> | null = null

// Status configuration
const statusConfig = computed(() => {
  switch (props.instance?.status) {
    case 'connected':
      return {
        label: t('evolution.status.connected'),
        color: 'text-success',
        bgColor: 'bg-success/10',
        icon: Check,
        description: t('evolution.connectModal.connectedDesc')
      }
    case 'connecting':
      return {
        label: t('evolution.status.connecting'),
        color: 'text-warning',
        bgColor: 'bg-warning/10',
        icon: Loader2,
        description: t('evolution.connectModal.connectingDesc')
      }
    case 'disconnected':
    default:
      return {
        label: t('evolution.status.disconnected'),
        color: 'text-muted-foreground',
        bgColor: 'bg-muted',
        icon: WifiOff,
        description: t('evolution.connectModal.disconnectedDesc')
      }
  }
})

const showQRCode = computed(() => {
  return props.instance?.status === 'connecting' || 
         (props.instance?.status === 'disconnected' && props.qrCode)
})

const showConnectedState = computed(() => {
  return props.instance?.status === 'connected'
})

// Poll for status updates
const startPolling = () => {
  if (pollingInterval) return
  
  isPolling.value = true
  pollingInterval = setInterval(() => {
    emit('refresh-qr')
  }, 3000) // Poll every 3 seconds
}

const stopPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
  isPolling.value = false
}

// Watch for connection success
watch(() => props.instance?.status, (status) => {
  if (status === 'connected') {
    stopPolling()
  }
})

// Watch for open state
watch(() => props.open, (isOpen) => {
  if (isOpen) {
    startPolling()
  } else {
    stopPolling()
  }
})

const handleClose = () => {
  stopPolling()
  emit('update:open', false)
}

const handleRefreshQR = () => {
  emit('refresh-qr')
}

const handleDisconnect = () => {
  emit('disconnect')
}

// Lifecycle
onUnmounted(() => {
  stopPolling()
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open && instance"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-background/80 backdrop-blur-sm"
          @click="handleClose"
        />

        <!-- Modal Content -->
        <div class="relative w-full max-w-md bg-card rounded-2xl border border-border shadow-2xl animate-scale-in">
          <!-- Header -->
          <div class="flex items-center justify-between p-4 border-b border-border">
            <div class="flex items-center gap-3">
              <div :class="['w-10 h-10 rounded-xl flex items-center justify-center', statusConfig.bgColor]">
                <component
                  :is="statusConfig.icon"
                  :class="['w-5 h-5', statusConfig.color, instance.status === 'connecting' && 'animate-spin']"
                />
              </div>
              <div>
                <h2 class="text-lg font-semibold text-foreground">
                  {{ instance.display_name || instance.instance_name }}
                </h2>
                <p class="text-sm text-muted-foreground">
                  {{ statusConfig.label }}
                </p>
              </div>
            </div>
            
            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8"
              @click="handleClose"
            >
              <X class="w-5 h-5" />
            </Button>
          </div>

          <!-- Content -->
          <div class="p-6">
            <!-- Connected State -->
            <div v-if="showConnectedState" class="text-center py-8 animate-fade-in">
              <div class="w-20 h-20 rounded-full bg-success/10 flex items-center justify-center mx-auto mb-4">
                <Check class="w-10 h-10 text-success" />
              </div>
              <h3 class="text-lg font-semibold text-foreground mb-2">
                {{ t('evolution.connectModal.connectedTitle') }}
              </h3>
              <p class="text-sm text-muted-foreground mb-2">
                {{ statusConfig.description }}
              </p>
              
              <!-- Connection Details -->
              <div v-if="instance.phone || instance.profile_name" class="mt-4 p-4 rounded-xl bg-muted/50 text-left">
                <div v-if="instance.phone" class="flex items-center gap-2 text-sm mb-2">
                  <Smartphone class="w-4 h-4 text-muted-foreground" />
                  <span class="text-muted-foreground">{{ t('evolution.connectModal.phone') }}:</span>
                  <span class="font-medium">{{ instance.phone }}</span>
                </div>
                <div v-if="instance.profile_name" class="flex items-center gap-2 text-sm">
                  <User class="w-4 h-4 text-muted-foreground" />
                  <span class="text-muted-foreground">{{ t('evolution.connectModal.profile') }}:</span>
                  <span class="font-medium">{{ instance.profile_name }}</span>
                </div>
              </div>

              <Button
                variant="outline"
                class="mt-6 gap-2"
                @click="handleDisconnect"
              >
                <WifiOff class="w-4 h-4" />
                {{ t('evolution.actions.disconnect') }}
              </Button>
            </div>

            <!-- QR Code State -->
            <div v-else-if="showQRCode">
              <QRCodeDisplay
                :qr-code="qrCode"
                :loading="loading"
                :error="error"
                :auto-refresh="true"
                @refresh="handleRefreshQR"
              />
            </div>

            <!-- Loading State (initial) -->
            <div v-else class="text-center py-8">
              <Loader2 class="w-12 h-12 text-primary animate-spin mx-auto mb-4" />
              <p class="text-sm text-muted-foreground">
                {{ t('evolution.connectModal.generatingQR') }}
              </p>
            </div>
          </div>

          <!-- Footer -->
          <div class="p-4 border-t border-border bg-muted/30 rounded-b-2xl">
            <div class="flex items-center gap-2 text-xs text-muted-foreground">
              <AlertTriangle class="w-4 h-4 flex-shrink-0" />
              <p>{{ t('evolution.connectModal.warning') }}</p>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
/* Modal animations */
.modal-enter-active {
  transition: opacity 0.3s ease;
}

.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

/* Scale animation */
.animate-scale-in {
  animation: scaleIn 0.35s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.92) translateY(8px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

/* Fade animation */
.animate-fade-in {
  animation: fadeIn 0.3s ease-out forwards;
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
  .modal-enter-active,
  .modal-leave-active,
  .animate-scale-in,
  .animate-fade-in {
    animation: none !important;
    transition: none !important;
  }
}
</style>
