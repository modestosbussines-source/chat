// OMNI Design System - Inbox Types

export interface Contact {
  id: string
  name: string | null
  phone: string
  avatar?: string
  isOnline?: boolean
  company?: string | null
  email?: string | null
  location?: string | null
  tags: Tag[]
  createdAt: string
  notes?: string | null
}

export interface Tag {
  id: string
  name: string
  color?: string
}

export interface Message {
  id: string
  content: string
  createdAt: string
  senderType: 'agent' | 'customer' | 'system'
  senderName?: string
  status?: 'sending' | 'sent' | 'delivered' | 'read' | 'failed'
  media?: MediaAttachment | null
  replyTo?: MessageReply | null
  reactions?: Reaction[]
}

export interface MediaAttachment {
  type: 'image' | 'audio' | 'video' | 'file' | 'location'
  url: string
  name?: string
  mimeType?: string
  size?: number
  duration?: number
}

export interface MessageReply {
  id: string
  content: string
  senderName: string
}

export interface Reaction {
  id: string
  emoji: string
  count: number
  userIds?: string[]
}

export interface Conversation {
  id: string
  contact: Contact | null
  lastMessage: Message | null
  unreadCount: number
  status: 'open' | 'pending' | 'resolved'
  channel: 'whatsapp' | 'instagram' | 'facebook' | 'telegram' | 'email'
  tags: Tag[]
  assignedTo: {
    id: string
    name: string
  } | null
  isStarred?: boolean
  createdAt?: string
  updatedAt?: string
}

export interface InboxFilter {
  value: 'all' | 'unread' | 'assigned' | 'unassigned'
  label: string
  count: number
}

export interface ConversationListProps {
  conversations: Conversation[]
  activeConversationId?: string
  loading?: boolean
  filter?: InboxFilter['value']
}

export interface MessageThreadProps {
  messages: Message[]
  contactName?: string
  contactAvatar?: string
  loading?: boolean
}

export interface MessageInputProps {
  placeholder?: string
  disabled?: boolean
  replyTo?: Message | null
  showEmojiPicker?: boolean
}
