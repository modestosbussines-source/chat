// OMNI Design System - Evolution API Composable
import { ref, computed, onUnmounted } from 'vue'
import { api } from '@/services/api'
import { wsService } from '@/services/websocket'
import { toast } from 'vue-sonner'
import type {
  EvolutionInstance,
  CreateInstanceRequest,
  QRCodeResponse,
  ConnectionStatusResponse,
  EvolutionInstanceStatus,
  InstanceFilters,
  QRCodeDisplayState
} from '@/types/evolution'

// State
const instances = ref<EvolutionInstance[]>([])
const isLoading = ref(false)
const selectedInstance = ref<EvolutionInstance | null>(null)
const filters = ref<InstanceFilters>({
  status: 'all',
  search: ''
})

// QR Code state per instance
const qrCodeStates = ref<Map<string, QRCodeDisplayState>>(new Map())

// Auto-refresh intervals
const refreshIntervals = ref<Map<string, ReturnType<typeof setInterval>>>(new Map())

export function useEvolution() {
  // Computed
  const filteredInstances = computed(() => {
    let result = [...instances.value]

    // Filter by status
    if (filters.value.status !== 'all') {
      result = result.filter(i => i.status === filters.value.status)
    }

    // Filter by search
    if (filters.value.search.trim()) {
      const search = filters.value.search.toLowerCase()
      result = result.filter(
        i =>
          i.instance_name.toLowerCase().includes(search) ||
          i.display_name?.toLowerCase().includes(search) ||
          i.phone?.includes(search)
      )
    }

    return result
  })

  const connectedInstances = computed(() =>
    instances.value.filter(i => i.status === 'connected')
  )

  const disconnectedInstances = computed(() =>
    instances.value.filter(i => i.status === 'disconnected')
  )

  const totalInstances = computed(() => instances.value.length)

  // API Methods
  async function fetchInstances(): Promise<void> {
    isLoading.value = true
    try {
      const response = await api.get('/evolution/instances')
      instances.value = response.data.data || []
    } catch (error) {
      console.error('Failed to fetch instances:', error)
      toast.error('Erro ao carregar instâncias')
      throw error
    } finally {
      isLoading.value = false
    }
  }

  async function createInstance(data: CreateInstanceRequest): Promise<EvolutionInstance> {
    try {
      const response = await api.post('/evolution/instances', data)
      const newInstance = response.data.data
      instances.value.unshift(newInstance)
      toast.success(`Instância "${newInstance.display_name || newInstance.instance_name}" criada!`)
      return newInstance
    } catch (error: any) {
      const message = error.response?.data?.message || 'Erro ao criar instância'
      toast.error(message)
      throw error
    }
  }

  async function deleteInstance(instanceId: string): Promise<void> {
    try {
      await api.delete(`/evolution/instances/${instanceId}`)
      instances.value = instances.value.filter(i => i.id !== instanceId)
      toast.success('Instância removida com sucesso')
    } catch (error: any) {
      const message = error.response?.data?.message || 'Erro ao remover instância'
      toast.error(message)
      throw error
    }
  }

  async function getInstance(instanceId: string): Promise<EvolutionInstance> {
    try {
      const response = await api.get(`/evolution/instances/${instanceId}`)
      return response.data.data
    } catch (error) {
      console.error('Failed to get instance:', error)
      throw error
    }
  }

  async function getQRCode(instanceId: string): Promise<QRCodeResponse> {
    try {
      const response = await api.get(`/evolution/instances/${instanceId}/qrcode`)
      const data = response.data.data

      // Update QR code state
      qrCodeStates.value.set(instanceId, {
        isLoading: false,
        qrCode: data.qrcode,
        expiresAt: Date.now() + 20000, // QR expires in ~20 seconds
        error: null
      })

      // Update instance status
      updateInstanceStatus(instanceId, data.status)

      return data
    } catch (error: any) {
      qrCodeStates.value.set(instanceId, {
        isLoading: false,
        qrCode: null,
        expiresAt: null,
        error: error.response?.data?.message || 'Erro ao gerar QR Code'
      })
      throw error
    }
  }

  async function getConnectionStatus(instanceId: string): Promise<ConnectionStatusResponse> {
    try {
      const response = await api.get(`/evolution/instances/${instanceId}/status`)
      return response.data.data
    } catch (error) {
      console.error('Failed to get connection status:', error)
      throw error
    }
  }

  async function disconnectInstance(instanceId: string): Promise<void> {
    try {
      await api.post(`/evolution/instances/${instanceId}/disconnect`)
      
      // Update local state
      updateInstanceStatus(instanceId, 'disconnected')
      
      // Clear QR code state
      qrCodeStates.value.delete(instanceId)
      
      // Stop auto-refresh
      stopQRCodeRefresh(instanceId)
      
      toast.success('Instância desconectada')
    } catch (error: any) {
      const message = error.response?.data?.message || 'Erro ao desconectar instância'
      toast.error(message)
      throw error
    }
  }

  // QR Code Auto-refresh
  function startQRCodeRefresh(instanceId: string, callback?: () => void): void {
    // Clear existing interval
    stopQRCodeRefresh(instanceId)

    // Fetch immediately
    fetchQRCodeWithState(instanceId)

    // Set up auto-refresh every 15 seconds
    const interval = setInterval(() => {
      fetchQRCodeWithState(instanceId)
      callback?.()
    }, 15000)

    refreshIntervals.value.set(instanceId, interval)
  }

  function stopQRCodeRefresh(instanceId: string): void {
    const interval = refreshIntervals.value.get(instanceId)
    if (interval) {
      clearInterval(interval)
      refreshIntervals.value.delete(instanceId)
    }
  }

  async function fetchQRCodeWithState(instanceId: string): Promise<void> {
    // Set loading state
    qrCodeStates.value.set(instanceId, {
      isLoading: true,
      qrCode: qrCodeStates.value.get(instanceId)?.qrCode || null,
      expiresAt: qrCodeStates.value.get(instanceId)?.expiresAt || null,
      error: null
    })

    try {
      await getQRCode(instanceId)
    } catch {
      // Error already handled in getQRCode
    }
  }

  // Instance state helpers
  function updateInstanceStatus(instanceId: string, status: EvolutionInstanceStatus): void {
    const instance = instances.value.find(i => i.id === instanceId)
    if (instance) {
      instance.status = status
      if (status === 'connected') {
        instance.connected_at = new Date().toISOString()
      } else if (status === 'disconnected') {
        instance.connected_at = null
        instance.phone = null
        instance.profile_name = null
      }
    }
  }

  function getQRCodeState(instanceId: string): QRCodeDisplayState {
    return (
      qrCodeStates.value.get(instanceId) || {
        isLoading: false,
        qrCode: null,
        expiresAt: null,
        error: null
      }
    )
  }

  // WebSocket setup
  function setupWebSocketListeners(): void {
    wsService.on('evolution.connection.update', (data: any) => {
      const { instance, state } = data
      if (instance && state) {
        const status = state === 'open' ? 'connected' : 'disconnected'
        updateInstanceStatus(instance, status)

        if (status === 'connected') {
          toast.success(`WhatsApp conectado na instância "${instance}"`)
          stopQRCodeRefresh(instance)
        } else if (status === 'disconnected') {
          toast.warning(`WhatsApp desconectado na instância "${instance}"`)
        }
      }
    })

    wsService.on('evolution.qrcode', (data: any) => {
      const { instance, qrcode } = data
      if (instance && qrcode) {
        qrCodeStates.value.set(instance, {
          isLoading: false,
          qrCode: qrcode,
          expiresAt: Date.now() + 20000,
          error: null
        })
      }
    })
  }

  // Cleanup
  function cleanup(): void {
    refreshIntervals.value.forEach(interval => clearInterval(interval))
    refreshIntervals.value.clear()
    qrCodeStates.value.clear()
  }

  onUnmounted(cleanup)

  return {
    // State
    instances,
    isLoading,
    selectedInstance,
    filters,

    // Computed
    filteredInstances,
    connectedInstances,
    disconnectedInstances,
    totalInstances,

    // Methods
    fetchInstances,
    createInstance,
    deleteInstance,
    getInstance,
    getQRCode,
    getConnectionStatus,
    disconnectInstance,

    // QR Code
    startQRCodeRefresh,
    stopQRCodeRefresh,
    fetchQRCodeWithState,
    getQRCodeState,

    // Helpers
    updateInstanceStatus,
    setupWebSocketListeners,
    cleanup
  }
}
