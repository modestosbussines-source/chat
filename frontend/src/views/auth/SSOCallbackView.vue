<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/services/api'
import { toast } from 'vue-sonner'
import { Loader2, AlertCircle, CheckCircle, ArrowLeft } from 'lucide-vue-next'

// OMNI Design System Components
import AuthLayout from '@/components/auth/AuthLayout.vue'
import { Button } from '@/components/ui/button'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()

const status = ref<'loading' | 'success' | 'error'>('loading')
const errorMessage = ref('')

onMounted(async () => {
  try {
    const response = await api.get('/me')
    const user = response.data.data

    authStore.setAuth({ user })

    status.value = 'success'
    toast.success(t('auth.ssoLoginSuccess'))

    setTimeout(() => {
      if (user.role?.name === 'agent') {
        router.push('/analytics/agents')
      } else {
        router.push('/')
      }
    }, 1000)
  } catch (error: any) {
    status.value = 'error'
    errorMessage.value = error.response?.data?.message || t('auth.ssoLoginFailed')
  }
})
</script>

<template>
  <AuthLayout title="">
    <div class="text-center py-8">
      <!-- Loading State -->
      <div v-if="status === 'loading'" class="space-y-4">
        <div class="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center mx-auto">
          <Loader2 class="w-10 h-10 text-primary animate-spin" />
        </div>
        <h2 class="text-xl font-semibold text-foreground">
          {{ $t('auth.ssoLoading') }}
        </h2>
        <p class="text-muted-foreground">
          {{ $t('auth.ssoLoadingDesc') }}
        </p>
      </div>

      <!-- Success State -->
      <div v-else-if="status === 'success'" class="space-y-4">
        <div class="w-20 h-20 rounded-full bg-success/10 flex items-center justify-center mx-auto">
          <CheckCircle class="w-10 h-10 text-success" />
        </div>
        <h2 class="text-xl font-semibold text-foreground">
          {{ $t('auth.ssoSuccess') }}
        </h2>
        <p class="text-muted-foreground">
          {{ $t('auth.ssoSuccessDesc') }}
        </p>
        <div class="flex items-center justify-center gap-2 text-sm text-muted-foreground">
          <Loader2 class="w-4 h-4 animate-spin" />
          Redirecionando...
        </div>
      </div>

      <!-- Error State -->
      <div v-else class="space-y-4">
        <div class="w-20 h-20 rounded-full bg-destructive/10 flex items-center justify-center mx-auto">
          <AlertCircle class="w-10 h-10 text-destructive" />
        </div>
        <h2 class="text-xl font-semibold text-foreground">
          {{ $t('auth.ssoFailed') }}
        </h2>
        <p class="text-muted-foreground">
          {{ errorMessage }}
        </p>
        <RouterLink to="/login">
          <Button variant="outline" class="mt-4">
            <ArrowLeft class="mr-2 w-4 h-4" />
            {{ $t('auth.returnToLogin') }}
          </Button>
        </RouterLink>
      </div>
    </div>
  </AuthLayout>
</template>
