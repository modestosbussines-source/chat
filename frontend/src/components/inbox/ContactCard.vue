<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Avatar from '@/components/ui/avatar/Avatar.vue'
import AvatarFallback from '@/components/ui/avatar/AvatarFallback.vue'
import AvatarImage from '@/components/ui/avatar/AvatarImage.vue'
import Badge from '@/components/ui/badge/Badge.vue'
import { Mail, Phone, MapPin, Calendar, Tag, Building2, User } from 'lucide-vue-next'
import type { Contact } from '@/types/inbox'

interface Props {
  contact: Contact
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compact: false
})

const emit = defineEmits<{
  (e: 'edit'): void
  (e: 'message'): void
  (e: 'call'): void
}>()

const { t } = useI18n()

const contactInitials = computed(() => {
  const name = props.contact?.name || ''
  if (!name) return '?'
  return name
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
})

const formattedPhone = computed(() => {
  if (!props.contact?.phone) return ''
  const phone = props.contact.phone.replace(/\D/g, '')
  if (phone.length === 13 && phone.startsWith('55')) {
    // Brazilian format: +55 (11) 99999-9999
    return `+${phone.slice(0, 2)} (${phone.slice(2, 4)}) ${phone.slice(4, 9)}-${phone.slice(9)}`
  }
  return props.contact.phone
})

const createdDate = computed(() => {
  if (!props.contact?.createdAt) return ''
  return new Date(props.contact.createdAt).toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
})
</script>

<template>
  <div :class="['contact-card', compact ? 'p-3' : 'p-4']">
    <!-- Header -->
    <div class="flex items-start gap-3 mb-4">
      <div class="relative flex-shrink-0">
        <Avatar :class="compact ? 'w-10 h-10' : 'w-14 h-14'">
          <AvatarImage :src="contact.avatar" :alt="contact.name" />
          <AvatarFallback :class="compact ? 'text-sm' : 'text-lg'" class="bg-primary/10 text-primary">
            {{ contactInitials }}
          </AvatarFallback>
        </Avatar>
        
        <!-- Online Indicator -->
        <div
          v-if="contact.isOnline"
          class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-success border-2 border-card"
        />
      </div>
      
      <div class="flex-1 min-w-0">
        <h3 class="font-semibold text-foreground truncate">
          {{ contact.name || 'Unknown Contact' }}
        </h3>
        <p v-if="contact.company" class="text-sm text-muted-foreground flex items-center gap-1 mt-0.5">
          <Building2 class="w-3 h-3" />
          {{ contact.company }}
        </p>
      </div>
      
      <!-- Actions -->
      <div class="flex gap-1">
        <button
          @click="emit('message')"
          class="p-2 rounded-lg hover:bg-muted transition-colors"
          :title="t('contacts.actions.message')"
        >
          <svg
            class="w-4 h-4 text-muted-foreground"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
          </svg>
        </button>
        <button
          v-if="contact.phone"
          @click="emit('call')"
          class="p-2 rounded-lg hover:bg-muted transition-colors"
          :title="t('contacts.actions.call')"
        >
          <Phone class="w-4 h-4 text-muted-foreground" />
        </button>
      </div>
    </div>
    
    <!-- Contact Info -->
    <div v-if="!compact" class="space-y-3">
      <!-- Phone -->
      <div v-if="contact.phone" class="flex items-center gap-3 text-sm">
        <div class="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
          <Phone class="w-4 h-4 text-primary" />
        </div>
        <div>
          <p class="text-xs text-muted-foreground">{{ t('contacts.fields.phone') }}</p>
          <p class="text-foreground">{{ formattedPhone }}</p>
        </div>
      </div>
      
      <!-- Email -->
      <div v-if="contact.email" class="flex items-center gap-3 text-sm">
        <div class="w-8 h-8 rounded-lg bg-secondary/10 flex items-center justify-center">
          <Mail class="w-4 h-4 text-secondary" />
        </div>
        <div>
          <p class="text-xs text-muted-foreground">{{ t('contacts.fields.email') }}</p>
          <p class="text-foreground truncate">{{ contact.email }}</p>
        </div>
      </div>
      
      <!-- Location -->
      <div v-if="contact.location" class="flex items-center gap-3 text-sm">
        <div class="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center">
          <MapPin class="w-4 h-4 text-muted-foreground" />
        </div>
        <div>
          <p class="text-xs text-muted-foreground">{{ t('contacts.fields.location') }}</p>
          <p class="text-foreground">{{ contact.location }}</p>
        </div>
      </div>
      
      <!-- Created Date -->
      <div class="flex items-center gap-3 text-sm">
        <div class="w-8 h-8 rounded-lg bg-muted flex items-center justify-center">
          <Calendar class="w-4 h-4 text-muted-foreground" />
        </div>
        <div>
          <p class="text-xs text-muted-foreground">{{ t('contacts.fields.created') }}</p>
          <p class="text-foreground">{{ createdDate }}</p>
        </div>
      </div>
    </div>
    
    <!-- Tags -->
    <div v-if="contact.tags && contact.tags.length > 0" class="mt-4">
      <div class="flex items-center gap-2 mb-2">
        <Tag class="w-3 h-3 text-muted-foreground" />
        <span class="text-xs text-muted-foreground">{{ t('contacts.fields.tags') }}</span>
      </div>
      <div class="flex flex-wrap gap-1.5">
        <Badge
          v-for="tag in contact.tags"
          :key="tag.id"
          variant="secondary"
          class="text-xs"
        >
          {{ tag.name }}
        </Badge>
      </div>
    </div>
    
    <!-- Notes -->
    <div v-if="contact.notes && !compact" class="mt-4 p-3 rounded-lg bg-muted">
      <p class="text-xs text-muted-foreground mb-1">{{ t('contacts.fields.notes') }}</p>
      <p class="text-sm text-foreground">{{ contact.notes }}</p>
    </div>
  </div>
</template>

<style scoped>
.contact-card {
  @apply rounded-xl bg-card border border-border;
}
</style>
