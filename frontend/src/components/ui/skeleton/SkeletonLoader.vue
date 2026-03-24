<script setup lang="ts">
interface Props {
  variant?: 'text' | 'avatar' | 'card' | 'conversation' | 'message'
  lines?: number
  animated?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'text',
  lines: 3,
  animated: true
})
</script>

<template>
  <!-- Text Skeleton -->
  <div v-if="variant === 'text'" class="space-y-2">
    <div
      v-for="i in lines"
      :key="i"
      :class="[
        'skeleton h-4 rounded-md',
        animated && 'skeleton-animated',
        i === lines && 'w-3/4'
      ]"
    />
  </div>

  <!-- Avatar Skeleton -->
  <div v-else-if="variant === 'avatar'" class="flex items-center gap-3">
    <div :class="['skeleton skeleton-avatar', animated && 'skeleton-animated']" />
    <div class="flex-1 space-y-2">
      <div :class="['skeleton h-4 w-24 rounded', animated && 'skeleton-animated']" />
      <div :class="['skeleton h-3 w-32 rounded', animated && 'skeleton-animated']" />
    </div>
  </div>

  <!-- Card Skeleton -->
  <div v-else-if="variant === 'card'" class="skeleton-card space-y-4">
    <div :class="['skeleton h-40 rounded-lg', animated && 'skeleton-animated']" />
    <div class="space-y-2">
      <div :class="['skeleton h-4 w-3/4 rounded', animated && 'skeleton-animated']" />
      <div :class="['skeleton h-4 w-1/2 rounded', animated && 'skeleton-animated']" />
    </div>
  </div>

  <!-- Conversation Item Skeleton -->
  <div v-else-if="variant === 'conversation'" class="space-y-3 p-3">
    <div v-for="i in lines" :key="i" class="flex items-center gap-3">
      <div :class="['skeleton w-12 h-12 rounded-full flex-shrink-0', animated && 'skeleton-animated']" />
      <div class="flex-1 space-y-2">
        <div class="flex justify-between">
          <div :class="['skeleton h-4 w-24 rounded', animated && 'skeleton-animated']" />
          <div :class="['skeleton h-3 w-12 rounded', animated && 'skeleton-animated']" />
        </div>
        <div :class="['skeleton h-3 w-full rounded', animated && 'skeleton-animated']" />
      </div>
    </div>
  </div>

  <!-- Message Skeleton -->
  <div v-else-if="variant === 'message'" class="space-y-4">
    <div v-for="i in lines" :key="i" :class="['flex', i % 2 === 0 ? 'justify-end' : 'justify-start']">
      <div
        :class="[
          'skeleton rounded-2xl',
          i % 2 === 0 ? 'w-48' : 'w-56',
          animated && 'skeleton-animated'
        ]"
        style="height: 48px;"
      />
    </div>
  </div>
</template>

<style scoped>
.skeleton-animated {
  position: relative;
  overflow: hidden;
  background: hsl(var(--muted));
}

.skeleton-animated::after {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  background: linear-gradient(
    90deg,
    transparent,
    hsl(var(--muted-foreground) / 0.08),
    transparent
  );
  animation: shimmer 1.5s infinite;
  transform: translateX(-100%);
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.skeleton-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: hsl(var(--muted));
}

.skeleton-card {
  padding: 1rem;
  border-radius: var(--radius);
  background: hsl(var(--card));
  border: 1px solid hsl(var(--border));
}
</style>
