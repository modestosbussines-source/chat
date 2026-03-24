<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ConversationItem from './ConversationItem.vue'
import SkeletonLoader from '@/components/ui/skeleton/SkeletonLoader.vue'
import type { Conversation } from '@/types/inbox'

interface Props {
  conversations: Conversation[]
  activeConversationId?: string
  loading?: boolean
  filter?: 'all' | 'unread' | 'assigned' | 'unassigned'
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  filter: 'all'
})

const emit = defineEmits<{
  (e: 'select', conversation: Conversation): void
  (e: 'filter-change', filter: string): void
}>()

const { t } = useI18n()

const searchQuery = ref('')
const isRefreshing = ref(false)

const filteredConversations = computed(() => {
  let result = [...props.conversations]

  // Apply filter
  switch (props.filter) {
    case 'unread':
      result = result.filter(c => c.unreadCount > 0)
      break
    case 'assigned':
      result = result.filter(c => c.assignedTo)
      break
    case 'unassigned':
      result = result.filter(c => !c.assignedTo)
      break
  }

  // Apply search
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(c =>
      c.contact?.name?.toLowerCase().includes(query) ||
      c.contact?.phone?.includes(query) ||
      c.lastMessage?.content?.toLowerCase().includes(query)
    )
  }

  // Sort by last message time (newest first)
  return result.sort((a, b) => {
    const dateA = new Date(a.lastMessage?.createdAt || 0).getTime()
    const dateB = new Date(b.lastMessage?.createdAt || 0).getTime()
    return dateB - dateA
  })
})

const filterOptions = computed(() => [
  { value: 'all', label: t('inbox.filters.all'), count: props.conversations.length },
  { value: 'unread', label: t('inbox.filters.unread'), count: props.conversations.filter(c => c.unreadCount > 0).length },
  { value: 'assigned', label: t('inbox.filters.assigned'), count: props.conversations.filter(c => c.assignedTo).length },
  { value: 'unassigned', label: t('inbox.filters.unassigned'), count: props.conversations.filter(c => !c.assignedTo).length }
])

const handleSelect = (conversation: Conversation) => {
  emit('select', conversation)
}

const handleFilterChange = (filter: string) => {
  emit('filter-change', filter)
}

const handleRefresh = async () => {
  isRefreshing.value = true
  // Animation delay
  setTimeout(() => {
    isRefreshing.value = false
  }, 1000)
}
</script>

<template>
  <div class="conversation-list flex flex-col h-full bg-card overflow-hidden">
    <!-- Header -->
    <div class="p-3 md:p-4 border-b border-border space-y-3">
      <h2 class="text-lg font-semibold text-foreground animate-fade-in">{{ t('inbox.title') }}</h2>
      
      <!-- Search -->
      <div class="relative animate-fade-in-up stagger-1">
        <svg
          class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground transition-colors duration-200"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="11" cy="11" r="8" />
          <path d="m21 21-4.3-4.3" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          :placeholder="t('inbox.search.placeholder')"
          class="w-full pl-10 pr-4 py-2.5 text-sm bg-muted border border-transparent rounded-xl focus:outline-none focus:border-ring focus:bg-background focus:ring-2 focus:ring-primary/10 transition-all duration-200"
        />
      </div>

      <!-- Filter Tabs -->
      <div class="flex gap-1 p-1 bg-muted rounded-xl animate-fade-in-up stagger-2">
        <button
          v-for="option in filterOptions"
          :key="option.value"
          @click="handleFilterChange(option.value)"
          :class="[
            'flex-1 flex items-center justify-center gap-1.5 px-2 py-2 text-xs font-medium rounded-lg transition-all duration-200 relative',
            filter === option.value
              ? 'bg-background text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'
          ]"
        >
          {{ option.label }}
          <span
            v-if="option.count > 0"
            :class="[
              'px-1.5 py-0.5 text-[10px] rounded-full transition-all duration-200',
              filter === option.value
                ? 'bg-primary text-primary-foreground'
                : 'bg-muted-foreground/15'
            ]"
          >
            {{ option.count }}
          </span>
        </button>
      </div>
    </div>

    <!-- Conversation List -->
    <div class="flex-1 overflow-y-auto">
      <!-- Loading State -->
      <div v-if="loading" class="p-3 space-y-3">
        <SkeletonLoader variant="conversation" :lines="6" />
      </div>

      <!-- Empty State -->
      <div
        v-else-if="filteredConversations.length === 0"
        class="flex flex-col items-center justify-center h-full p-6 text-center"
      >
        <div class="w-16 h-16 mb-4 rounded-full bg-muted flex items-center justify-center animate-float">
          <svg
            class="w-8 h-8 text-muted-foreground"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
          </svg>
        </div>
        <p class="text-sm font-medium text-foreground animate-fade-in-up stagger-1">
          {{ t('inbox.empty.title') }}
        </p>
        <p class="text-xs text-muted-foreground mt-1 animate-fade-in-up stagger-2">
          {{ t('inbox.empty.description') }}
        </p>
      </div>

      <!-- Conversation Items -->
      <div v-else class="py-2">
        <TransitionGroup name="list" tag="div">
          <ConversationItem
            v-for="(conversation, index) in filteredConversations"
            :key="conversation.id"
            :conversation="conversation"
            :active="conversation.id === activeConversationId"
            :class="['animate-fade-in-up', `stagger-${Math.min(index + 1, 8)}`]"
            @click="handleSelect(conversation)"
          />
        </TransitionGroup>
      </div>
    </div>

    <!-- Refresh indicator -->
    <Transition name="fade">
      <div
        v-if="isRefreshing"
        class="absolute inset-0 bg-background/50 backdrop-blur-sm flex items-center justify-center z-10"
      >
        <div class="flex items-center gap-2 px-4 py-2 bg-card rounded-full shadow-lg">
          <svg class="w-5 h-5 animate-spin text-primary" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <span class="text-sm text-muted-foreground">Atualizando...</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
/* List animations */
.list-enter-active {
  animation: listItemEnter 0.35s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.list-leave-active {
  animation: listItemLeave 0.2s ease-in forwards;
}

.list-move {
  transition: transform 0.3s ease;
}

@keyframes listItemEnter {
  from {
    opacity: 0;
    transform: translateX(-12px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateX(0) scale(1);
  }
}

@keyframes listItemLeave {
  to {
    opacity: 0;
    transform: translateX(-12px) scale(0.98);
  }
}

/* Fade transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Animation utilities */
.animate-fade-in {
  animation: fadeIn 0.3s ease-out forwards;
}

.animate-fade-in-up {
  animation: fadeInUp 0.4s ease-out forwards;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

/* Stagger delays */
.stagger-1 { animation-delay: 0.05s; }
.stagger-2 { animation-delay: 0.1s; }
.stagger-3 { animation-delay: 0.15s; }
.stagger-4 { animation-delay: 0.2s; }
.stagger-5 { animation-delay: 0.25s; }
.stagger-6 { animation-delay: 0.3s; }
.stagger-7 { animation-delay: 0.35s; }
.stagger-8 { animation-delay: 0.4s; }

/* Touch device optimizations */
@media (hover: none) {
  .conversation-item:active {
    transform: scale(0.98);
  }
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .list-enter-active,
  .list-leave-active,
  .animate-fade-in,
  .animate-fade-in-up,
  .animate-float {
    animation: none !important;
    transition: none !important;
  }
}
</style>
