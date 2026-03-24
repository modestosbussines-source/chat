<script setup lang="ts">
interface Props {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  text?: string
  fullScreen?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  fullScreen: false
})

const sizeClasses = {
  sm: 'w-4 h-4 border-2',
  md: 'w-8 h-8 border-2',
  lg: 'w-12 h-12 border-3',
  xl: 'w-16 h-16 border-4'
}
</script>

<template>
  <div
    :class="[
      'flex flex-col items-center justify-center gap-4',
      fullScreen && 'min-h-screen bg-background'
    ]"
  >
    <!-- Spinner -->
    <div
      :class="[
        sizeClasses[props.size],
        'rounded-full border-primary/20 border-t-primary animate-spin'
      ]"
    />
    
    <!-- Text -->
    <p v-if="text" class="text-sm text-muted-foreground animate-pulse">
      {{ text }}
    </p>
  </div>
</template>
