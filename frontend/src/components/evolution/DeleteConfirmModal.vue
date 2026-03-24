<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { AlertTriangle, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'

interface Props {
  open: boolean
  instanceName?: string
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  instanceName: '',
  loading: false
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'confirm'): void
}>()

const { t } = useI18n()

const handleConfirm = () => {
  emit('confirm')
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <div class="flex items-center gap-3 mb-2">
          <div class="w-12 h-12 rounded-full bg-destructive/10 flex items-center justify-center flex-shrink-0">
            <AlertTriangle class="w-6 h-6 text-destructive" />
          </div>
          <div>
            <DialogTitle class="text-lg">
              {{ t('evolution.deleteModal.title') }}
            </DialogTitle>
            <DialogDescription class="mt-1">
              {{ t('evolution.deleteModal.description') }}
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div class="py-4">
        <div class="p-4 rounded-xl bg-muted">
          <p class="text-sm text-foreground font-medium">
            {{ instanceName }}
          </p>
          <p class="text-xs text-muted-foreground mt-1">
            {{ t('evolution.deleteModal.warning') }}
          </p>
        </div>
      </div>

      <DialogFooter class="gap-2 sm:gap-0">
        <Button
          variant="outline"
          @click="emit('update:open', false)"
          :disabled="loading"
        >
          {{ t('common.cancel') }}
        </Button>
        <Button
          variant="destructive"
          @click="handleConfirm"
          :disabled="loading"
        >
          <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
          {{ t('evolution.deleteModal.confirm') }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
