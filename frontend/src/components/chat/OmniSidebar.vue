<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { omnisService } from '@/services/api'
import { toast } from 'vue-sonner'
import { Loader2, Play, Image as ImageIcon, FileAudio, FileVideo, FileText, ChevronDown, ChevronRight, X, Settings, UploadCloud } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog'

const props = defineProps<{
  contactId: string
}>()

const emit = defineEmits(['close'])
const { t } = useI18n()

interface OmniScript {
  id: string
  title: string
  media_type: string
}

interface OmniCategory {
  id: string
  name: string
  color: string
  Scripts: OmniScript[]
}

const categories = ref<OmniCategory[]>([])
const isLoading = ref(true)
const sendingScriptId = ref<string | null>(null)
const expandedCategories = ref<Set<string>>(new Set())

const isSettingsOpen = ref(false)
const isUploadingZip = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)

async function fetchCategories() {
  try {
    const res = await omnisService.listCategories()
    categories.value = res.data.data || res.data
    if (categories.value.length > 0) {
      expandedCategories.value.add(categories.value[0].id)
    }
  } catch (error) {
    toast.error('Failed to load Omnis')
  } finally {
    isLoading.value = false
  }
}

function toggleCategory(id: string) {
  if (expandedCategories.value.has(id)) {
    expandedCategories.value.delete(id)
  } else {
    expandedCategories.value.add(id)
  }
}

async function sendScript(scriptId: string) {
  if (sendingScriptId.value) return
  sendingScriptId.value = scriptId
  try {
    await omnisService.send(props.contactId, scriptId)
    toast.success('Omni disparado com sucesso!')
  } catch (e: any) {
    const msg = e.response?.data?.message || 'Falha ao disparar Omni'
    toast.error(msg)
  } finally {
    sendingScriptId.value = null
  }
}

function getIconForMedia(type: string) {
  switch (type) {
    case 'image': return ImageIcon
    case 'audio': return FileAudio
    case 'voice': return FileAudio
    case 'video': return FileVideo
    case 'document': return FileText
    default: return FileText
  }
}

onMounted(() => {
  fetchCategories()
})

function handleZipSelect(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  if (!file.name.endsWith('.zip')) {
    toast.error('Somente arquivos .zip são aceitos.')
    return
  }

  isUploadingZip.value = true
  omnisService.importZip(file)
    .then(() => {
      toast.success('Omnis importados com sucesso!')
      isSettingsOpen.value = false
      isLoading.value = true
      fetchCategories()
    })
    .catch((e: any) => {
      const msg = e.response?.data?.message || 'Falha ao importar Omnis'
      toast.error(msg)
    })
    .finally(() => {
      isUploadingZip.value = false
      input.value = '' // reset
    })
}
</script>

<template>
  <div class="flex flex-col h-full bg-white dark:bg-[#1a1a1a] border-l border-neutral-200 dark:border-neutral-800">
    <div class="flex items-center justify-between p-4 border-b border-neutral-200 dark:border-neutral-800">
      <h3 class="font-semibold text-neutral-900 dark:text-neutral-100 flex items-center gap-2">
        <Play class="w-4 h-4" />
        Omnis
      </h3>
      <div class="flex items-center gap-1">
        <Button variant="ghost" size="icon" @click="isSettingsOpen = true" class="h-8 w-8 text-neutral-500 hover:text-neutral-900 dark:hover:text-white">
          <Settings class="w-4 h-4" />
        </Button>
        <Button variant="ghost" size="icon" @click="emit('close')" class="h-8 w-8 text-neutral-500 hover:text-neutral-900 dark:hover:text-white">
          <X class="w-4 h-4" />
        </Button>
      </div>
    </div>

    <div v-if="isLoading" class="flex-1 flex items-center justify-center">
      <Loader2 class="w-6 h-6 animate-spin text-neutral-400" />
    </div>

    <ScrollArea v-else class="flex-1">
      <div v-if="categories.length === 0" class="p-4 text-center text-sm text-neutral-500">
        Nenhum script configurado.
      </div>
      <div v-for="category in categories" :key="category.id" class="border-b border-neutral-100 dark:border-neutral-800/50 last:border-0">
        <button 
          @click="toggleCategory(category.id)"
          class="w-full flex items-center justify-between p-3 hover:bg-neutral-50 dark:hover:bg-neutral-800/50 transition-colors text-left"
        >
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 rounded-full" :style="{ backgroundColor: category.color || '#000' }"></div>
            <span class="text-sm font-medium text-neutral-900 dark:text-neutral-100">{{ category.name }}</span>
            <span class="text-xs text-neutral-400 bg-neutral-100 dark:bg-neutral-800 px-1.5 py-0.5 rounded-full">
              {{ category.Scripts?.length || 0 }}
            </span>
          </div>
          <ChevronDown v-if="expandedCategories.has(category.id)" class="w-4 h-4 text-neutral-400" />
          <ChevronRight v-else class="w-4 h-4 text-neutral-400" />
        </button>
        
        <div v-show="expandedCategories.has(category.id)" class="p-2 pt-0 space-y-1 bg-neutral-50/50 dark:bg-neutral-900/20">
          <button
            v-for="script in category.Scripts"
            :key="script.id"
            @click="sendScript(script.id)"
            :disabled="sendingScriptId === script.id"
            class="w-full flex items-center gap-3 p-2 rounded-md hover:bg-neutral-200 dark:hover:bg-neutral-800 transition-colors text-left group disabled:opacity-50"
          >
            <div class="p-1.5 rounded bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-700 text-neutral-500 dark:text-neutral-400 group-hover:text-black dark:group-hover:text-white transition-colors">
              <Loader2 v-if="sendingScriptId === script.id" class="w-3.5 h-3.5 animate-spin" />
              <component v-else :is="getIconForMedia(script.media_type)" class="w-3.5 h-3.5" />
            </div>
            <span class="text-sm text-neutral-700 dark:text-neutral-300 group-hover:text-neutral-900 dark:group-hover:text-white truncate flex-1 leading-tight">
              {{ script.title }}
            </span>
          </button>
          <div v-if="!category.Scripts || category.Scripts.length === 0" class="p-2 text-xs text-neutral-400 text-center">
            Vazio
          </div>
        </div>
      </div>
    </ScrollArea>

    <Dialog v-model:open="isSettingsOpen">
      <DialogContent class="max-w-sm">
        <DialogHeader>
          <DialogTitle>Configurações Omnis</DialogTitle>
          <DialogDescription>
            Faça upload do seu arquivo de funil (.zip contendo o config.json e mídias) para alimentar as respostas globais do Omni.
          </DialogDescription>
        </DialogHeader>
        <div class="py-4 flex flex-col items-center justify-center gap-4">
          <input
            ref="fileInputRef"
            type="file"
            accept=".zip"
            class="hidden"
            @change="handleZipSelect"
          />
          <Button 
            @click="fileInputRef?.click()" 
            :disabled="isUploadingZip"
            class="w-full h-24 border-2 border-dashed flex-col gap-2"
            variant="outline"
          >
            <Loader2 v-if="isUploadingZip" class="w-6 h-6 animate-spin text-neutral-500" />
            <UploadCloud v-else class="w-6 h-6 text-neutral-500" />
            <span class="text-sm font-medium text-neutral-600">
              {{ isUploadingZip ? 'Importando e Processando...' : 'Selecionar Arquivo .ZIP' }}
            </span>
          </Button>
          <p class="text-xs text-neutral-400 text-center px-4">
            Aviso: O upload destruirá os Omnis anteriores. O arquivo deve conter um "config.json" na raiz.
          </p>
        </div>
      </DialogContent>
    </Dialog>
  </div>
</template>
