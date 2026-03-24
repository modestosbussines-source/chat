<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  MessageSquare,
  ChevronLeft,
  ChevronRight,
  Menu,
  X,
  Home,
  Inbox,
  Users,
  Settings,
  LogOut
} from 'lucide-vue-next'
import { wsService } from '@/services/websocket'
import { authService } from '@/services/api'
import OrganizationSwitcher from './OrganizationSwitcher.vue'
import UserMenu from './UserMenu.vue'
import ActiveCallPanel from '@/components/calling/ActiveCallPanel.vue'
import { navigationItems } from './navigation'
import OmniLogo from '@/components/brand/OmniLogo.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const isCollapsed = ref(false)
const isMobileMenuOpen = ref(false)
const isDragging = ref(false)
const dragStartX = ref(0)

onMounted(() => {
  if (authStore.isAuthenticated) {
    authStore.refreshUserData()

    wsService.connect(async () => {
      try {
        const resp = await authService.getWSToken()
        return resp.data.data.token
      } catch {
        return null
      }
    })
  }
})

// Close mobile menu on route change
watch(() => route.path, () => {
  isMobileMenuOpen.value = false
})

const navigation = computed(() => {
  return navigationItems
    .filter(item => {
      if (item.childPermissions) {
        return item.childPermissions.some(p => authStore.hasPermission(p, 'read'))
      }
      return !item.permission || authStore.hasPermission(item.permission, 'read')
    })
    .map(item => {
      const filteredChildren = item.children?.filter(
        child => !child.permission || authStore.hasPermission(child.permission, 'read')
      )

      let effectivePath = item.path
      if (item.childPermissions && item.permission && !authStore.hasPermission(item.permission, 'read') && filteredChildren?.length) {
        effectivePath = filteredChildren[0].path
      }

      const originalPath = item.path
      const isActive = originalPath === '/'
        ? route.name === 'dashboard'
        : originalPath === '/chat'
          ? route.name === 'chat' || route.name === 'chat-conversation'
          : originalPath === '/inbox'
            ? route.name === 'inbox'
            : route.path.startsWith(originalPath)

      return {
        ...item,
        path: effectivePath,
        active: isActive,
        children: filteredChildren
      }
    })
})

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

// Swipe gesture handling for mobile
const handleTouchStart = (e: TouchEvent) => {
  dragStartX.value = e.touches[0].clientX
}

const handleTouchMove = (e: TouchEvent) => {
  const currentX = e.touches[0].clientX
  const diff = currentX - dragStartX.value
  
  if (Math.abs(diff) > 50) {
    if (diff > 0 && !isMobileMenuOpen.value) {
      isMobileMenuOpen.value = true
    } else if (diff < 0 && isMobileMenuOpen.value) {
      isMobileMenuOpen.value = false
    }
  }
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div 
    class="app-layout flex h-screen bg-background overflow-hidden"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
  >
    <!-- Skip link for accessibility -->
    <a href="#main-content" class="skip-link">{{ $t('nav.skipToMain') }}</a>

    <!-- Mobile Overlay -->
    <Transition name="fade">
      <div
        v-if="isMobileMenuOpen"
        class="fixed inset-0 z-40 bg-background/80 backdrop-blur-sm md:hidden"
        @click="closeMobileMenu"
      />
    </Transition>

    <!-- Mobile Header -->
    <header class="fixed top-0 left-0 right-0 z-30 flex h-14 items-center justify-between border-b border-border bg-background/95 backdrop-blur-md px-4 md:hidden">
      <RouterLink to="/" class="flex items-center gap-2" @click="closeMobileMenu">
        <OmniLogo size="sm" :show-text="false" />
        <span class="font-semibold text-foreground">OMNI</span>
      </RouterLink>
      
      <div class="flex items-center gap-2">
        <!-- Quick Inbox Button -->
        <RouterLink to="/inbox">
          <Button variant="ghost" size="icon" class="relative h-9 w-9">
            <Inbox class="h-5 w-5" />
          </Button>
        </RouterLink>
        
        <!-- Menu Toggle -->
        <Button
          variant="ghost"
          size="icon"
          class="h-9 w-9"
          aria-label="Toggle menu"
          :aria-expanded="isMobileMenuOpen"
          @click="toggleMobileMenu"
        >
          <Transition name="scale" mode="out-in">
            <X v-if="isMobileMenuOpen" class="h-5 w-5" />
            <Menu v-else class="h-5 w-5" />
          </Transition>
        </Button>
      </div>
    </header>

    <!-- Mobile Navigation Drawer -->
    <Transition name="slide-left">
      <aside
        v-if="isMobileMenuOpen"
        class="fixed inset-y-0 left-0 z-50 w-72 bg-card border-r border-border shadow-xl md:hidden"
      >
        <div class="flex flex-col h-full">
          <!-- Drawer Header -->
          <div class="flex items-center justify-between p-4 border-b border-border">
            <OmniLogo size="md" />
            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8"
              @click="closeMobileMenu"
            >
              <X class="h-5 w-5" />
            </Button>
          </div>

          <!-- Navigation Items -->
          <ScrollArea class="flex-1 py-4">
            <nav class="space-y-1 px-3">
              <RouterLink
                v-for="(item, index) in navigation"
                :key="item.path"
                :to="item.path"
                :class="[
                  'flex items-center gap-3 rounded-lg px-3 py-3 text-sm font-medium transition-all duration-200 animate-fade-in-up',
                  item.active
                    ? 'bg-primary text-primary-foreground shadow-sm'
                    : 'text-muted-foreground hover:text-foreground hover:bg-muted',
                  `stagger-${index + 1}`
                ]"
                @click="closeMobileMenu"
              >
                <component :is="item.icon" class="h-5 w-5 shrink-0" />
                <span>{{ $t(item.name) }}</span>
                <span v-if="item.children" class="ml-auto">
                  <ChevronRight class="h-4 w-4" />
                </span>
              </RouterLink>
            </nav>
          </ScrollArea>

          <!-- Drawer Footer -->
          <div class="border-t border-border p-4">
            <UserMenu @logout="handleLogout" />
          </div>
        </div>
      </aside>
    </Transition>

    <!-- Desktop Sidebar -->
    <aside
      :class="[
        'hidden md:flex flex-col border-r border-border bg-card transition-all duration-300',
        'relative',
        isCollapsed ? 'w-16 lg:w-20' : 'w-64 lg:w-72'
      ]"
      role="navigation"
      aria-label="Main navigation"
    >
      <!-- Logo -->
      <div class="flex h-14 lg:h-16 items-center justify-between px-4 border-b border-border">
        <RouterLink to="/" class="flex items-center gap-2">
          <Transition name="scale" mode="out-in">
            <OmniLogo v-if="isCollapsed" key="small" size="sm" :show-text="false" />
            <OmniLogo v-else key="large" size="md" />
          </Transition>
        </RouterLink>
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8 opacity-0 lg:opacity-100 transition-opacity"
          :aria-label="isCollapsed ? $t('nav.expandSidebar') : $t('nav.collapseSidebar')"
          :aria-expanded="!isCollapsed"
          @click="toggleSidebar"
        >
          <Transition name="rotate" mode="out-in">
            <ChevronRight v-if="isCollapsed" key="expand" class="h-4 w-4" />
            <ChevronLeft v-else key="collapse" class="h-4 w-4" />
          </Transition>
        </Button>
      </div>

      <!-- Navigation -->
      <ScrollArea class="flex-1 py-4">
        <nav class="space-y-1 px-3" role="menubar">
          <RouterLink
            v-for="item in navigation"
            :key="item.path"
            :to="item.path"
            :class="[
              'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200 group relative',
              item.active
                ? 'bg-primary text-primary-foreground shadow-sm shadow-primary/20'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted',
              isCollapsed && 'lg:justify-center lg:px-2'
            ]"
            role="menuitem"
            :aria-current="item.active ? 'page' : undefined"
          >
            <component :is="item.icon" class="h-5 w-5 shrink-0" aria-hidden="true" />
            <Transition name="fade">
              <span v-if="!isCollapsed" class="lg:block">{{ $t(item.name) }}</span>
            </Transition>
            
            <!-- Tooltip for collapsed state -->
            <div
              v-if="isCollapsed"
              class="absolute left-full ml-2 px-2 py-1 bg-popover text-popover-foreground text-sm rounded-md shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap z-50"
            >
              {{ $t(item.name) }}
            </div>

            <!-- Active indicator -->
            <div
              v-if="item.active"
              class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-6 bg-primary rounded-r-full"
            />
          </RouterLink>
        </nav>
      </ScrollArea>

      <!-- User Menu -->
      <div class="border-t border-border">
        <UserMenu :collapsed="isCollapsed" @logout="handleLogout" />
      </div>
    </aside>

    <!-- Main Content -->
    <main 
      id="main-content" 
      class="flex-1 overflow-hidden pt-14 md:pt-0" 
      role="main"
    >
      <RouterView v-slot="{ Component, route: viewRoute }">
        <Transition name="page" mode="out-in">
          <component :is="Component" :key="viewRoute.path" />
        </Transition>
      </RouterView>
      <ActiveCallPanel />
    </main>
  </div>
</template>

<style scoped>
/* Mobile drawer animation */
.slide-left-enter-active,
.slide-left-leave-active {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-left-enter-from,
.slide-left-leave-to {
  transform: translateX(-100%);
}

/* Page transition */
.page-enter-active {
  animation: page-enter 0.3s ease-out;
}

.page-leave-active {
  animation: page-leave 0.2s ease-in;
}

@keyframes page-enter {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes page-leave {
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(-8px);
  }
}

/* Scale transition */
.scale-enter-active,
.scale-leave-active {
  transition: transform 0.15s ease-out, opacity 0.15s ease-out;
}

.scale-enter-from,
.scale-leave-to {
  opacity: 0;
  transform: scale(0.8);
}

/* Rotate transition */
.rotate-enter-active,
.rotate-leave-active {
  transition: transform 0.2s ease-out;
}

.rotate-enter-from {
  transform: rotate(-90deg);
}

.rotate-leave-to {
  transform: rotate(90deg);
}

/* Fade transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Touch device optimizations */
@media (hover: none) {
  .group:hover .group-hover\:opacity-100 {
    opacity: 0;
  }
  
  .group:active .group-hover\:opacity-100 {
    opacity: 1;
  }
}

/* Large screens - auto collapse */
@media (min-width: 1536px) {
  .app-layout {
    --sidebar-collapsed-width: 5rem;
  }
}
</style>
