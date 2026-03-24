<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import InstanceList from '@/components/evolution/InstanceList.vue'
import ConnectModal from '@/components/evolution/ConnectModal.vue'
import CreateInstanceModal from '@/components/evolution/CreateInstanceModal.vue'
import DeleteConfirmModal from '@/components/evolution/DeleteConfirmModal.vue'
import { useEvolution } from '@/composables/useEvolution'
import type { EvolutionInstance, EvolutionInstanceStatus } from '@/types/evolution'

const { t } = useI18n()

const {
  instances,
  isLoading,
  filters,
  fetchInstances,
  createInstance,
  deleteInstance,
  getQRCode,
  disconnectInstance,
  startQRCodeRefresh,
  stopQRCodeRefresh,
  setupWebSocketListeners
} = useEvolution()

// Modal states
const isConnectModalOpen = ref(false)
const isCreateModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const selectedInstance = ref<EvolutionInstance | null>(null)
const isCreating = ref(false)
const isDeleting = ref(false)
const qrCode = ref<string | null>(null)
const qrError = ref<string | null>(null)

// Handlers
const handleCreate = async (data: { instance_name: string; display_name: string }) => {
  isCreating.value = true
  try {
    const newInstance = await createInstance(data)
    isCreateModalOpen.value = false
    
    // Auto-open connect modal for new instance
    selectedInstance.value = newInstance
    await handleConnect(newInstance)
  } catch (error) {
    // Error already handled in composable
  } finally {
    isCreating.value = false
  }
}

const handleConnect = async (instance: EvolutionInstance) => {
  selectedInstance.value = instance
  qrCode.value = null
  qrError.value = null
  isConnectModalOpen.value = true
  
  try {
    const response = await getQRCode(instance.id)
    qrCode.value = response.qrcode
    
    // Start auto-refresh for QR code
    startQRCodeRefresh(instance.id, async () => {
      try {
        const response = await getQRCode(instance.id)
        qrCode.value = response.qrcode
        qrError.value = null
      } catch (error: any) {
        qrError.value = error.response?.data?.message || 'Erro ao atualizar QR Code'
      }
    })
  } catch (error: any) {
    qrError.value = error.response?.data?.message || 'Erro ao gerar QR Code'
  }
}

const handleDisconnect = async (instance: EvolutionInstance) => {
  try {
    await disconnectInstance(instance.id)
    stopQRCodeRefresh(instance.id)
    
    // Update selected instance if it's the same
    if (selectedInstance.value?.id === instance.id) {
      selectedInstance.value = {
        ...selectedInstance.value,
        status: 'disconnected',
        phone: null,
        profile_name: null,
        connected_at: null
      }
    }
  } catch (error) {
    // Error already handled in composable
  }
}

const handleDeleteClick = (instance: EvolutionInstance) => {
  selectedInstance.value = instance
  isDeleteModalOpen.value = true
}

const handleDeleteConfirm = async () => {
  if (!selectedInstance.value) return
  
  isDeleting.value = true
  try {
    await deleteInstance(selectedInstance.value.id)
    isDeleteModalOpen.value = false
    selectedInstance.value = null
    
    // Close connect modal if it was for this instance
    if (isConnectModalOpen.value && selectedInstance.value?.id) {
      isConnectModalOpen.value = false
    }
  } catch (error) {
    // Error already handled in composable
  } finally {
    isDeleting.value = false
  }
}

const handleViewDetails = (instance: EvolutionInstance) => {
  // TODO: Implement details view
  toast.info(`${t('evolution.detailsView')}: ${instance.display_name || instance.instance_name}`)
}

const handleRefreshQR = async () => {
  if (!selectedInstance.value) return
  
  try {
    const response = await getQRCode(selectedInstance.value.id)
    qrCode.value = response.qrcode
    qrError.value = null
  } catch (error: any) {
    qrError.value = error.response?.data?.message || 'Erro ao atualizar QR Code'
  }
}

const handleFilterChange = (filter: EvolutionInstanceStatus | 'all') => {
  filters.value.status = filter
}

// Lifecycle
onMounted(async () => {
  await fetchInstances()
  setupWebSocketListeners()
})

onUnmounted(() => {
  // Cleanup is handled by composable
})
</script>

<template>
  <div class="evolution-view h-full">
    <InstanceList
      :instances="instances"
      :loading="isLoading"
      @create="isCreateModalOpen = true"
      @connect="handleConnect"
      @disconnect="handleDisconnect"
      @delete="handleDeleteClick"
      @view-details="handleViewDetails"
      @filter-change="handleFilterChange"
    />

    <!-- Connect Modal (QR Code) -->
    <ConnectModal
      v-model:open="isConnectModalOpen"
      :instance="selectedInstance"
      :qr-code="qrCode"
      :error="qrError"
      :loading="isLoading"
      @refresh-qr="handleRefreshQR"
      @disconnect="handleDisconnect(selectedInstance!)"
    />

    <!-- Create Instance Modal -->
    <CreateInstanceModal
      v-model:open="isCreateModalOpen"
      :loading="isCreating"
      @create="handleCreate"
    />

    <!-- Delete Confirmation Modal -->
    <DeleteConfirmModal
      v-model:open="isDeleteModalOpen"
      :instance-name="selectedInstance?.display_name || selectedInstance?.instance_name || ''"
      :loading="isDeleting"
      @confirm="handleDeleteConfirm"
    />
  </div>
</template>

<style scoped>
.evolution-view {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
