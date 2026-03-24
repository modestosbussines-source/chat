<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plus, Search, Filter, Wifi, WifiOff, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import InstanceCard from './InstanceCard.vue'
import SkeletonLoader from '@/components/ui/skeleton/SkeletonLoader.vue'
import type { EvolutionInstance, EvolutionInstanceStatus } from '@/types/evolution'

interface Props {
  instances: EvolutionInstance[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  (e: 'create'): void
  (e: 'connect', instance: EvolutionInstance): void
  (e: 'disconnect', instance: EvolutionInstance): void
  (e: 'delete', instance: EvolutionInstance): void
  (e: 'view-details', instance: EvolutionInstance): void
  (e: 'filter-change', filter: EvolutionInstanceStatus | 'all'): void
}>()

const { t } = useI18n()

const searchQuery = ref('')
const activeFilter = ref<EvolutionInstanceStatus | 'all'>('all')

const filterOptions = computed(() => [
  { 
    value: 'all' as const, 
    label: t('evolution.filters.all'),
    count: props.instances.length,
    icon: null
  },
  { 
    value: 'connected' as const, 
    label: t('evolution.filters.connected'),
    count: props.instances.filter(i => i.status === 'connected').length,
    icon: Wifi
  },
  { 
    value: 'connecting' as const, 
    label: t('evolution.filters.connecting'),
    count: props.instances.filter(i => i.status === 'connecting').length,
    icon: Loader2
  },
  { 
    value: 'disconnected' as const, 
    label: t('evolution.filters.disconnected'),
    count: props.instances.filter(i => i.status === 'disconnected').length,
    icon: WifiOff
  }
])

const filteredInstances = computed(() => {
  let result = [...props.instances]

  // Apply status filter
  if (activeFilter.value !== 'all') {
    result = result.filter(i => i.status === activeFilter.value)
  }

  // Apply search
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(
      i =>
        i.instance_name.toLowerCase().includes(query) ||
        i.display_name?.toLowerCase().includes(query) ||
        i.phone?.includes(query)
    )
  }

  return result
})

const handleFilterChange = (filter: EvolutionInstanceStatus | 'all') => {
  activeFilter.value = filter
  emit('filter-change', filter)
}
</script>

<template>
  <div class="instance-list flex flex-col h-full">
    <!-- Header -->
    <div class="p-4 border-b border-border space-y-4">
      <!-- Title & Create Button -->
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold text-foreground">{{ t('evolution.title') }}</h2>
          <p class="text-sm text-muted-foreground">{{ t('evolution.subtitle') }}</p>
        </div>
        <Button size="sm" class="gap-2" @click="emit('create')">
          <Plus class="w-4 h-4" />
          <span class="hidden sm:inline">{{ t('evolution.actions.create') }}</span>
        </Button>
      </div>

      <!-- Search -->
      <div class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <input
          v-model="searchQuery"
          type="text"
          :placeholder="t('evolution.searchPlaceholder')"
          class="w-full pl-10 pr-4 py-2.5 text-sm bg-muted border border-transparent rounded-xl focus:outline-none focus:border-ring focus:bg-background focus:ring-2 focus:ring-primary/10 transition-all duration-200"
        />
      </div>

      <!-- Filter Tabs -->
      <div class="flex gap-1 p-1 bg-muted rounded-xl overflow-x-auto">
        <button
          v-for="option in filterOptions"
          :key="option.value"
          @click="handleFilterChange(option.value)"
          :class="[
            'flex items-center justify-center gap-1.5 px-3 py-2 text-xs font-medium rounded-lg transition-all duration-200 whitespace-nowrap',
            activeFilter === option.value
              ? 'bg-background text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'
          ]"
        >
          <component v-if="option.icon" :is="option.icon" class="w-3.5 h-3.5" />
          {{ option.label }}
          <span
            :class="[
              'px-1.5 py-0.5 text-[10px] rounded-full',
              activeFilter === option.value
                ? 'bg-primary text-primary-foreground'
                : 'bg-muted-foreground/15'
            ]"
          >
            {{ option.count }}
          </span>
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto p-4">
      <!-- Loading State -->
      <div v-if="loading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 6" :key="i" class="rounded-xl border border-border p-4">
          <div class="flex items-center gap-3 mb-4">
            <div class="skeleton w-12 h-12 rounded-full" />
            <div class="flex-1 space-y-2">
              <div class="skeleton h-4 w-24 rounded" />
              <div class="skeleton h-3 w-32 rounded" />
            </div>
          </div>
          <div class="skeleton h-3 w-full rounded mb-2" />
          <div class="skeleton h-8 w-full rounded" />
        </div>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="filteredInstances.length === 0"
        class="flex flex-col items-center justify-center h-full text-center py-12"
      >
        <div class="w-20 h-20 mb-4 rounded-2xl bg-muted flex items-center justify-center animate-float">
          <svg
            class="w-10 h-10 text-muted-foreground"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path d="M17.5 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0Z" />
            <path d="M19.5 10.5a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z" />
            <path d="M21 21a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v16Z" />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-foreground mb-2">
          {{ searchQuery || activeFilter !== 'all' ? t('evolution.empty.searchTitle') : t('evolution.empty.title') }}
        </h3>
        <p class="text-sm text-muted-foreground mb-6 max-w-sm">
          {{ searchQuery || activeFilter !== 'all' 
            ? t('evolution.empty.searchDescription') 
            : t('evolution.empty.description') 
          }}
        </p>
        <Button
          v-if="!searchQuery && activeFilter === 'all'"
          class="gap-2"
          @click="emit('create')"
        >
          <Plus class="w-4 h-4" />
          {{ t('evolution.actions.createFirst') }}
        </Button>
      </div>

      <!-- Instance Grid -->
      <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <TransitionGroup name="list">
          <InstanceCard
            v-for="(instance, index) in filteredInstances"
            :key="instance.id"
            :instance="instance"
            class="animate-fade-in-up"
            :style="{ animationDelay: `${index * 50}ms` }"
            @connect="emit('connect', $event)"
            @disconnect="emit('disconnect', $event)"
            @delete="emit('delete', $event)"
            @view-details="emit('view-details', $event)"
          />
        </TransitionGroup>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* List animations */
.list-enter-active {
  animation: fadeInUp 0.4s ease-out forwards;
}

.list-leave-active {
  animation: fadeOut 0.2s ease-in forwards;
  position: absolute;
}

.list-move {
  transition: transform 0.3s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeOut {
  to {
    opacity: 0;
    transform: scale(0.95);
  }
}

/* Animation utilities */
.animate-fade-in-up {
  animation: fadeInUp 0.4s ease-out forwards;
  opacity: 0;
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

/* Skeleton loader */
.skeleton {
  background: hsl(var(--muted));
  position: relative;
  overflow: hidden;
}

.skeleton::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent, hsl(var(--muted-foreground) / 0.08), transparent);
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

/* Mobile optimizations */
@media (max-width: 640px) {
  .grid {
    grid-template-columns: 1fr;
  }
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .list-enter-active,
  .list-leave-active,
  .animate-fade-in-up,
  .animate-float {
    animation: none !important;
    opacity: 1 !important;
  }
}
</style>
