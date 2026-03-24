<script setup lang="ts">
import { useColorMode } from '@/composables/useColorMode'
import OmniLogo from '@/components/brand/OmniLogo.vue'

interface Props {
  title: string
  subtitle?: string
}

defineProps<Props>()

const { isDark, toggleTheme } = useColorMode()
</script>

<template>
  <div class="auth-layout min-h-screen flex relative overflow-hidden">
    <!-- Background Gradient -->
    <div class="absolute inset-0 bg-gradient-to-br from-background via-background to-primary/5 transition-all duration-700" />
    
    <!-- Decorative Blobs - Subtle Animation -->
    <div class="absolute top-0 right-0 w-full md:w-1/2 h-full overflow-hidden pointer-events-none">
      <div class="absolute top-20 right-10 md:right-20 w-64 md:w-96 h-64 md:h-96 bg-primary/8 rounded-full blur-3xl animate-float" />
      <div class="absolute bottom-20 left-10 w-48 md:w-72 h-48 md:h-72 bg-secondary/8 rounded-full blur-3xl animate-float" style="animation-delay: 1s;" />
    </div>
    
    <!-- Grid Pattern Overlay -->
    <div class="absolute inset-0 bg-[linear-gradient(rgba(120,119,198,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(120,119,198,0.03)_1px,transparent_1px)] bg-[size:64px_64px] pointer-events-none" />
    
    <!-- Theme Toggle -->
    <button
      @click="toggleTheme"
      class="absolute top-4 right-4 z-10 p-2.5 rounded-xl bg-card/80 backdrop-blur-sm border border-border hover:bg-card transition-all duration-200 hover:scale-105 active:scale-95"
      :title="isDark ? 'Modo claro' : 'Modo escuro'"
    >
      <svg
        v-if="isDark"
        class="w-5 h-5 text-foreground transition-transform duration-300 rotate-0"
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
        class="w-5 h-5 text-foreground transition-transform duration-300"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z" />
      </svg>
    </button>

    <!-- Left Side - Branding (hidden on mobile) -->
    <div class="hidden lg:flex lg:w-1/2 xl:w-2/5 relative items-center justify-center p-12">
      <div class="max-w-md">
        <!-- Logo with animation -->
        <div class="mb-8 animate-fade-in">
          <OmniLogo size="xl" />
        </div>
        
        <!-- Tagline -->
        <h2 class="text-3xl font-bold text-foreground mb-4 leading-tight animate-fade-in-up stagger-1">
          Plataforma omnichannel para
          <span class="text-primary bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">equipes de alta performance</span>
        </h2>
        
        <p class="text-muted-foreground mb-8 animate-fade-in-up stagger-2">
          Gerencie WhatsApp, Instagram, Facebook e muito mais em uma única inbox intuitiva.
          Feito para o mercado brasileiro.
        </p>

        <!-- Features with icons -->
        <div class="space-y-4">
          <div class="flex items-center gap-3 animate-fade-in-up stagger-3">
            <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center transition-transform hover:scale-110">
              <svg class="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <div>
              <h4 class="font-medium text-foreground">Inbox Unificada</h4>
              <p class="text-sm text-muted-foreground">Todas as conversas em um lugar</p>
            </div>
          </div>
          
          <div class="flex items-center gap-3 animate-fade-in-up stagger-4">
            <div class="w-10 h-10 rounded-xl bg-secondary/10 flex items-center justify-center transition-transform hover:scale-110">
              <svg class="w-5 h-5 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div>
              <h4 class="font-medium text-foreground">Automações Inteligentes</h4>
              <p class="text-sm text-muted-foreground">Chatbots e fluxos automatizados</p>
            </div>
          </div>
          
          <div class="flex items-center gap-3 animate-fade-in-up stagger-5">
            <div class="w-10 h-10 rounded-xl bg-coral/10 flex items-center justify-center transition-transform hover:scale-110">
              <svg class="w-5 h-5 text-coral" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div>
              <h4 class="font-medium text-foreground">Analytics Avançado</h4>
              <p class="text-sm text-muted-foreground">Métricas e insights em tempo real</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Side - Form -->
    <div class="flex-1 flex items-center justify-center p-4 sm:p-6 lg:p-12">
      <div class="w-full max-w-md">
        <!-- Mobile Logo -->
        <div class="lg:hidden flex justify-center mb-8 animate-bounce-in">
          <OmniLogo size="lg" />
        </div>

        <!-- Card -->
        <div class="auth-card rounded-2xl border border-border bg-card/80 backdrop-blur-xl p-6 sm:p-8 shadow-xl animate-scale-in">
          <!-- Header -->
          <div class="text-center mb-8">
            <h1 class="text-2xl font-bold text-foreground animate-fade-in-up">{{ title }}</h1>
            <p v-if="subtitle" class="text-muted-foreground mt-2 animate-fade-in-up stagger-1">{{ subtitle }}</p>
          </div>

          <!-- Content Slot -->
          <slot class="animate-fade-in-up stagger-2" />

          <!-- Footer -->
          <div class="mt-8 pt-6 border-t border-border animate-fade-in-up stagger-3">
            <p class="text-center text-sm text-muted-foreground">
              <slot name="footer">
                <!-- Default footer content -->
              </slot>
            </p>
          </div>
        </div>

        <!-- Copyright -->
        <p class="text-center text-xs text-muted-foreground mt-6 animate-fade-in">
          {{ new Date().getFullYear() }} OMNICall. Todos os direitos reservados.
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-layout {
  min-height: 100vh;
}

.auth-card {
  box-shadow: 
    0 0 0 1px hsl(var(--border) / 0.5),
    0 4px 6px -1px hsl(var(--primary) / 0.05),
    0 20px 25px -5px hsl(var(--primary) / 0.08);
}

/* Animation utilities */
.animate-fade-in {
  animation: fadeIn 0.4s ease-out forwards;
}

.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out forwards;
}

.animate-scale-in {
  animation: scaleIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.animate-bounce-in {
  animation: bounceIn 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55) forwards;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes bounceIn {
  0% {
    opacity: 0;
    transform: scale(0.3);
  }
  50% {
    transform: scale(1.05);
  }
  70% {
    transform: scale(0.95);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

/* Float animation */
.animate-float {
  animation: float 4s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-16px) rotate(2deg);
  }
}

/* Stagger delays */
.stagger-1 { animation-delay: 0.1s; }
.stagger-2 { animation-delay: 0.2s; }
.stagger-3 { animation-delay: 0.3s; }
.stagger-4 { animation-delay: 0.4s; }
.stagger-5 { animation-delay: 0.5s; }

/* Mobile optimizations */
@media (max-width: 640px) {
  .auth-card {
    @apply p-5;
    box-shadow: 
      0 0 0 1px hsl(var(--border) / 0.3),
      0 10px 15px -3px hsl(var(--primary) / 0.1);
  }
  
  .auth-layout {
    padding: env(safe-area-inset-top, 0) env(safe-area-inset-right, 0) env(safe-area-inset-bottom, 0) env(safe-area-inset-left, 0);
  }
}

/* Touch device optimizations */
@media (hover: none) {
  .auth-card:hover {
    transform: none;
  }
  
  button:active {
    transform: scale(0.95);
  }
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .animate-fade-in,
  .animate-fade-in-up,
  .animate-scale-in,
  .animate-bounce-in,
  .animate-float {
    animation: none !important;
    opacity: 1 !important;
    transform: none !important;
  }
  
  * {
    transition-duration: 0.01ms !important;
  }
}

/* Dark mode adjustments */
@media (prefers-color-scheme: dark) {
  .auth-card {
    box-shadow: 
      0 0 0 1px hsl(var(--border) / 0.3),
      0 4px 6px -1px hsl(var(--primary) / 0.15),
      0 20px 25px -5px hsl(var(--primary) / 0.2);
  }
}
</style>
