// Evolution API Types - OMNI Design System

export type EvolutionInstanceStatus = 'connected' | 'disconnected' | 'connecting'

export interface EvolutionInstance {
  id: string
  organization_id: string
  instance_name: string
  display_name: string
  phone: string | null
  status: EvolutionInstanceStatus
  profile_name: string | null
  profile_pic: string | null
  webhook_url: string
  last_error: string | null
  connected_at: string | null
  last_activity: string | null
  created_at: string
  updated_at: string
}

export interface CreateInstanceRequest {
  instance_name: string
  display_name: string
}

export interface CreateInstanceResponse {
  id: string
  instance_name: string
  display_name: string
  status: EvolutionInstanceStatus
  webhook_url: string
  created_at: string
}

export interface QRCodeResponse {
  qrcode: string // Base64 encoded QR code image
  status: EvolutionInstanceStatus
}

export interface ConnectionStatusResponse {
  status: EvolutionInstanceStatus
  phone: string | null
  profile_name: string | null
  connected_at: string | null
  evolution_data?: {
    state: string
    connection: string
    qrcode?: {
      base64: string
      timestamp: number
    }
  }
}

export interface InstanceMetrics {
  total_messages: number
  sent_today: number
  received_today: number
  last_message_at: string | null
}

// WebSocket Events
export interface EvolutionWebSocketEvent {
  event: 'connection.update' | 'qrcode' | 'messages.upsert' | 'instance.created'
  instance: string
  data: {
    state?: string
    qrcode?: string
    message?: unknown
  }
}

// UI State Types
export interface InstanceFilters {
  status: 'all' | EvolutionInstanceStatus
  search: string
}

export interface QRCodeDisplayState {
  isLoading: boolean
  qrCode: string | null
  expiresAt: number | null
  error: string | null
}
