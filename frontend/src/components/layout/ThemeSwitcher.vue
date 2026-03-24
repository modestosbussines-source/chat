<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Sun, Moon, Monitor } from 'lucide-vue-next'
import { useColorMode } from '@/composables/useColorMode'

const { colorMode, setColorMode } = useColorMode()
</script>

<template>
  <div class="flex gap-1 p-1 bg-muted rounded-lg" role="radiogroup" aria-label="Color theme">
    <Button
      variant="ghost"
      size="icon"
      class="h-8 w-8 rounded-md"
      :class="colorMode === 'light' && 'bg-background shadow-sm text-foreground'"
      :aria-checked="colorMode === 'light'"
      aria-label="Light theme"
      role="radio"
      @click="setColorMode('light')"
    >
      <Sun class="h-4 w-4" aria-hidden="true" />
    </Button>
    <Button
      variant="ghost"
      size="icon"
      class="h-8 w-8 rounded-md"
      :class="colorMode === 'dark' && 'bg-background shadow-sm text-foreground'"
      :aria-checked="colorMode === 'dark'"
      aria-label="Dark theme"
      role="radio"
      @click="setColorMode('dark')"
    >
      <Moon class="h-4 w-4" aria-hidden="true" />
    </Button>
    <Button
      variant="ghost"
      size="icon"
      class="h-8 w-8 rounded-md"
      :class="colorMode === 'system' && 'bg-background shadow-sm text-foreground'"
      :aria-checked="colorMode === 'system'"
      aria-label="System theme"
      role="radio"
      @click="setColorMode('system')"
    >
      <Monitor class="h-4 w-4" aria-hidden="true" />
    </Button>
  </div>
</template>
