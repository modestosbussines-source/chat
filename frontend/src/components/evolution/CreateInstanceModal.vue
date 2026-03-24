<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { X, Loader2, Smartphone, Tag } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { toast } from 'vue-sonner'

interface Props {
  open: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'create', data: { instance_name: string; display_name: string }): void
}>()

const { t } = useI18n()

const instanceName = ref('')
const displayName = ref('')

// Generate slug from display name
const generateSlug = (name: string): string => {
  return name
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '') // Remove diacritics
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

// Auto-generate instance name when display name changes
const handleDisplayNameChange = () => {
  instanceName.value = generateSlug(displayName.value)
}

const isValid = computed(() => {
  return instanceName.value.trim().length >= 3 && displayName.value.trim().length > 0
})

const handleClose = () => {
  instanceName.value = ''
  displayName.value = ''
  emit('update:open', false)
}

const handleSubmit = () => {
  if (!isValid.value) {
    toast.error(t('evolution.createModal.validation'))
    return
  }

  emit('create', {
    instance_name: instanceName.value.trim().toLowerCase(),
    display_name: displayName.value.trim()
  })
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
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
              <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center">
                <Smartphone class="w-5 h-5 text-primary" />
              </div>
              <div>
                <h2 class="text-lg font-semibold text-foreground">
                  {{ t('evolution.createModal.title') }}
                </h2>
                <p class="text-sm text-muted-foreground">
                  {{ t('evolution.createModal.subtitle') }}
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

          <!-- Form -->
          <form @submit.prevent="handleSubmit" class="p-6 space-y-5">
            <!-- Display Name -->
            <div class="space-y-2">
              <label for="displayName" class="text-sm font-medium text-foreground flex items-center gap-2">
                <Tag class="w-4 h-4 text-muted-foreground" />
                {{ t('evolution.createModal.displayName') }}
                <span class="text-destructive">*</span>
              </label>
              <input
                id="displayName"
                v-model="displayName"
                type="text"
                :placeholder="t('evolution.createModal.displayNamePlaceholder')"
                class="w-full px-4 py-3 text-sm bg-muted border border-transparent rounded-xl focus:outline-none focus:border-ring focus:bg-background focus:ring-2 focus:ring-primary/10 transition-all duration-200"
                @input="handleDisplayNameChange"
                required
              />
              <p class="text-xs text-muted-foreground">
                {{ t('evolution.createModal.displayNameHint') }}
              </p>
            </div>

            <!-- Instance Name (Slug) -->
            <div class="space-y-2">
              <label for="instanceName" class="text-sm font-medium text-foreground">
                {{ t('evolution.createModal.instanceName') }}
                <span class="text-destructive">*</span>
              </label>
              <div class="flex items-center">
                <span class="px-3 py-3 text-sm bg-muted border border-r-0 border-border rounded-l-xl text-muted-foreground">
                  omni-
                </span>
                <input
                  id="instanceName"
                  v-model="instanceName"
                  type="text"
                  :placeholder="t('evolution.createModal.instanceNamePlaceholder')"
                  class="flex-1 px-4 py-3 text-sm bg-muted border border-transparent rounded-r-xl focus:outline-none focus:border-ring focus:bg-background focus:ring-2 focus:ring-primary/10 transition-all duration-200 font-mono"
                  pattern="[a-z0-9-]+"
                  required
                />
              </div>
              <p class="text-xs text-muted-foreground">
                {{ t('evolution.createModal.instanceNameHint') }}
              </p>
            </div>

            <!-- Preview Card -->
            <div class="p-4 rounded-xl bg-muted/50 border border-border">
              <p class="text-xs text-muted-foreground mb-2">{{ t('evolution.createModal.preview') }}</p>
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-sm font-semibold text-primary">
                  {{ (displayName || 'OM').slice(0, 2).toUpperCase() }}
                </div>
                <div>
                  <p class="text-sm font-medium text-foreground">
                    {{ displayName || t('evolution.createModal.previewName') }}
                  </p>
                  <p class="text-xs text-muted-foreground font-mono">
                    omni-{{ instanceName || 'instance' }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex gap-3 pt-2">
              <Button
                type="button"
                variant="outline"
                class="flex-1"
                @click="handleClose"
              >
                {{ t('common.cancel') }}
              </Button>
              <Button
                type="submit"
                class="flex-1"
                :disabled="!isValid || loading"
              >
                <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
                {{ t('evolution.createModal.create') }}
              </Button>
            </div>
          </form>
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
  .animate-scale-in {
    animation: none !important;
    transition: none !important;
  }
}
</style>
