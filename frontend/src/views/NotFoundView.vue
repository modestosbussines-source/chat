<script setup lang="ts">
import { useColorMode } from '@/composables/useColorMode'
import { Button } from '@/components/ui/button'
import { Home, ArrowLeft, Search } from 'lucide-vue-next'
import OmniLogo from '@/components/brand/OmniLogo.vue'

const { isDark, toggleTheme } = useColorMode()
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-background p-4 relative">
    <!-- Theme Toggle -->
    <button
      @click="toggleTheme"
      class="absolute top-4 right-4 p-2 rounded-lg hover:bg-muted transition-colors"
      :title="isDark ? 'Modo claro' : 'Modo escuro'"
    >
      <svg
        v-if="isDark"
        class="w-5 h-5 text-muted-foreground"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <circle cx="12" cy="12" r="4" />
        <path d="M12 2v2" /><path d="M12 20v2" />
        <path d="m4.93 4.93 1.41 1.41" /><path d="m17.66 17.66 1.41 1.41" />
        <path d="M2 12h2" /><path d="M20 12h2" />
        <path d="m6.34 17.66-1.41 1.41" /><path d="m19.07 4.93-1.41 1.41" />
      </svg>
      <svg
        v-else
        class="w-5 h-5 text-muted-foreground"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z" />
      </svg>
    </button>

    <!-- Background Decorations -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-1/4 left-1/4 w-64 h-64 bg-primary/5 rounded-full blur-3xl" />
      <div class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-secondary/5 rounded-full blur-3xl" />
    </div>

    <!-- Content -->
    <div class="relative text-center max-w-md">
      <!-- Logo -->
      <div class="mb-8">
        <OmniLogo size="md" />
      </div>

      <!-- 404 Number -->
      <div class="relative mb-6">
        <h1 class="text-[120px] font-bold text-muted-foreground/10 leading-none tracking-tight">
          404
        </h1>
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="w-24 h-24 rounded-2xl bg-muted/50 flex items-center justify-center">
            <Search class="w-12 h-12 text-muted-foreground" />
          </div>
        </div>
      </div>

      <!-- Text -->
      <h2 class="text-2xl font-bold text-foreground mb-2">
        {{ $t('notFound.title') }}
      </h2>
      <p class="text-muted-foreground mb-8">
        {{ $t('notFound.description') }}
      </p>

      <!-- Actions -->
      <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
        <Button variant="outline" class="w-full sm:w-auto" @click="$router.back()">
          <ArrowLeft class="h-4 w-4 mr-2" />
          {{ $t('notFound.goBack') }}
        </Button>
        <RouterLink to="/" class="w-full sm:w-auto">
          <Button class="w-full sm:w-auto">
            <Home class="h-4 w-4 mr-2" />
            {{ $t('notFound.goHome') }}
          </Button>
        </RouterLink>
      </div>

      <!-- Help Text -->
      <p class="text-sm text-muted-foreground mt-8">
        Precisa de ajuda?
        <a href="mailto:suporte@omnicalled.com" class="text-primary hover:underline">
          Entre em contato
        </a>
      </p>
    </div>

    <!-- Footer -->
    <p class="absolute bottom-4 text-xs text-muted-foreground">
      {{ new Date().getFullYear() }} OMNICall. Todos os direitos reservados.
    </p>
  </div>
</template>
