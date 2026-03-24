<script setup lang="ts">
interface Props {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  showText?: boolean
  variant?: 'default' | 'light' | 'dark'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  showText: true,
  variant: 'default'
})

const sizeClasses = {
  sm: { logo: 'w-8 h-8', text: 'text-lg' },
  md: { logo: 'w-10 h-10', text: 'text-xl' },
  lg: { logo: 'w-14 h-14', text: 'text-2xl' },
  xl: { logo: 'w-20 h-20', text: 'text-3xl' }
}

const currentSize = sizeClasses[props.size]
</script>

<template>
  <div class="flex items-center gap-3">
    <!-- Logo Mark -->
    <div
      :class="[
        currentSize.logo,
        'rounded-xl flex items-center justify-center',
        'bg-gradient-to-br from-primary to-primary-700',
        'shadow-lg shadow-primary/20'
      ]"
    >
      <svg
        :class="[props.size === 'sm' ? 'w-4 h-4' : props.size === 'md' ? 'w-5 h-5' : props.size === 'lg' ? 'w-7 h-7' : 'w-10 h-10']"
        viewBox="0 0 24 24"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <!-- OMNI Logo - Interconnected circles representing multi-channel -->
        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5" class="text-primary-foreground" />
        <circle cx="12" cy="8" r="3" stroke="currentColor" stroke-width="1.5" class="text-primary-foreground" />
        <circle cx="8" cy="15" r="3" stroke="currentColor" stroke-width="1.5" class="text-primary-foreground" />
        <circle cx="16" cy="15" r="3" stroke="currentColor" stroke-width="1.5" class="text-primary-foreground" />
        <line x1="12" y1="11" x2="10" y2="13" stroke="currentColor" stroke-width="1.5" class="text-primary-foreground" />
        <line x1="12" y1="11" x2="14" y2="13" stroke="currentColor" stroke-width="1.5" class="text-primary-foreground" />
        <line x1="10" y1="15" x2="14" y2="15" stroke="currentColor" stroke-width="1.5" class="text-primary-foreground" />
      </svg>
    </div>

    <!-- Logo Text -->
    <div v-if="showText" class="flex flex-col">
      <span :class="[currentSize.text, 'font-bold tracking-tight text-foreground']">
        OMNI
      </span>
      <span class="text-xs text-muted-foreground -mt-1 font-medium">
        Called
      </span>
    </div>
  </div>
</template>
