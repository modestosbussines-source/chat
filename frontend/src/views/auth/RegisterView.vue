<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { toast } from 'vue-sonner'
import { Loader2, User, Mail, Lock, ArrowRight, AlertCircle } from 'lucide-vue-next'

// OMNI Design System Components
import AuthLayout from '@/components/auth/AuthLayout.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const fullName = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const isLoading = ref(false)
const acceptTerms = ref(false)

const organizationId = computed(() => (route.query.org as string) || '')

// Password strength indicator
const passwordStrength = computed(() => {
  const pwd = password.value
  if (!pwd) return { score: 0, label: '', color: '' }
  
  let score = 0
  if (pwd.length >= 8) score++
  if (pwd.length >= 12) score++
  if (/[A-Z]/.test(pwd)) score++
  if (/[a-z]/.test(pwd)) score++
  if (/[0-9]/.test(pwd)) score++
  if (/[^A-Za-z0-9]/.test(pwd)) score++
  
  if (score <= 2) return { score, label: 'Fraca', color: 'bg-destructive' }
  if (score <= 4) return { score, label: 'Média', color: 'bg-warning' }
  return { score, label: 'Forte', color: 'bg-success' }
})

const passwordsMatch = computed(() => {
  if (!confirmPassword.value) return true
  return password.value === confirmPassword.value
})

const handleRegister = async () => {
  if (!organizationId.value) {
    toast.error(t('auth.invitationRequired'))
    return
  }

  if (!fullName.value || !email.value || !password.value) {
    toast.error(t('auth.fillAllFields'))
    return
  }

  if (password.value !== confirmPassword.value) {
    toast.error(t('auth.passwordsMismatch'))
    return
  }

  if (password.value.length < 8) {
    toast.error(t('auth.passwordTooShort'))
    return
  }

  if (!acceptTerms.value) {
    toast.error(t('auth.acceptTermsRequired'))
    return
  }

  isLoading.value = true

  try {
    await authStore.register({
      full_name: fullName.value,
      email: email.value,
      password: password.value,
      organization_id: organizationId.value
    })
    toast.success(t('auth.registrationSuccess'))
    router.push('/')
  } catch (error: any) {
    const message = error.response?.data?.message || t('auth.registrationFailed')
    toast.error(message)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <AuthLayout
    :title="$t('auth.createAccount')"
    :subtitle="$t('auth.createAccountDesc')"
  >
    <!-- No org ID - Show invitation required message -->
    <div v-if="!organizationId" class="text-center py-4">
      <div class="w-16 h-16 rounded-full bg-warning/10 flex items-center justify-center mx-auto mb-4">
        <AlertCircle class="w-8 h-8 text-warning" />
      </div>
      <p class="text-muted-foreground mb-6">
        {{ $t('auth.invitationRequired') }}
      </p>
      <RouterLink to="/login">
        <Button variant="outline" class="w-full">
          {{ $t('auth.signIn') }}
        </Button>
      </RouterLink>
    </div>

    <!-- Registration Form -->
    <form v-else @submit.prevent="handleRegister" class="space-y-5">
      <!-- Full Name -->
      <div class="space-y-2">
        <Label for="fullName" class="text-sm font-medium">
          {{ $t('auth.fullName') }}
        </Label>
        <div class="relative">
          <User class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            id="fullName"
            v-model="fullName"
            type="text"
            :placeholder="$t('auth.fullNamePlaceholder')"
            :disabled="isLoading"
            autocomplete="name"
            class="pl-10"
          />
        </div>
      </div>

      <!-- Email -->
      <div class="space-y-2">
        <Label for="email" class="text-sm font-medium">
          {{ $t('common.email') }}
        </Label>
        <div class="relative">
          <Mail class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            id="email"
            v-model="email"
            type="email"
            :placeholder="$t('auth.emailPlaceholder')"
            :disabled="isLoading"
            autocomplete="email"
            class="pl-10"
          />
        </div>
      </div>

      <!-- Password -->
      <div class="space-y-2">
        <Label for="password" class="text-sm font-medium">
          {{ $t('auth.password') }}
        </Label>
        <div class="relative">
          <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            id="password"
            v-model="password"
            type="password"
            :placeholder="$t('auth.passwordMinLength')"
            :disabled="isLoading"
            autocomplete="new-password"
            class="pl-10"
          />
        </div>
        <!-- Password Strength Indicator -->
        <div v-if="password" class="space-y-1">
          <div class="flex gap-1">
            <div
              v-for="i in 5"
              :key="i"
              :class="[
                'h-1 flex-1 rounded-full transition-colors',
                i <= passwordStrength.score ? passwordStrength.color : 'bg-muted'
              ]"
            />
          </div>
          <p class="text-xs text-muted-foreground">
            Força: <span :class="passwordStrength.color.replace('bg-', 'text-')">{{ passwordStrength.label }}</span>
          </p>
        </div>
      </div>

      <!-- Confirm Password -->
      <div class="space-y-2">
        <Label for="confirmPassword" class="text-sm font-medium">
          {{ $t('auth.confirmPassword') }}
        </Label>
        <div class="relative">
          <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            id="confirmPassword"
            v-model="confirmPassword"
            type="password"
            :placeholder="$t('auth.confirmPasswordPlaceholder')"
            :disabled="isLoading"
            autocomplete="new-password"
            class="pl-10"
            :class="{ 'border-destructive': !passwordsMatch }"
          />
        </div>
        <p v-if="!passwordsMatch" class="text-xs text-destructive">
          {{ $t('auth.passwordsMismatch') }}
        </p>
      </div>

      <!-- Terms Acceptance -->
      <div class="flex items-start gap-2">
        <input
          v-model="acceptTerms"
          type="checkbox"
          id="terms"
          class="mt-1 w-4 h-4 rounded border-border text-primary focus:ring-primary"
        />
        <label for="terms" class="text-sm text-muted-foreground cursor-pointer">
          Li e aceito os
          <a href="/terms" target="_blank" class="text-primary hover:underline">Termos de Uso</a>
          e
          <a href="/privacy" target="_blank" class="text-primary hover:underline">Política de Privacidade</a>
        </label>
      </div>

      <!-- Submit Button -->
      <Button
        type="submit"
        class="w-full h-11"
        :disabled="isLoading || !passwordsMatch"
      >
        <Loader2 v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
        <template v-else>
          {{ $t('auth.createAccountBtn') }}
          <ArrowRight class="ml-2 w-4 h-4" />
        </template>
      </Button>
    </form>

    <!-- Footer -->
    <template #footer>
      {{ $t('auth.alreadyHaveAccount') }}
      <RouterLink to="/login" class="text-primary font-medium hover:underline ml-1">
        {{ $t('auth.signIn') }}
      </RouterLink>
    </template>
  </AuthLayout>
</template>
