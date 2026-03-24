<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { PageHeader } from '@/components/shared'
import { api } from '@/services/api'
import { toast } from 'vue-sonner'
import { Plus, MessageSquare, Image, Video, FileIcon, Mic, Loader2, Pencil, Trash2, Upload, Send, Download, Import } from 'lucide-vue-next'

const { t } = useI18n()

interface OmniCategory {
  id: string
  name: string
  color: string
  scripts: OmniScript[]
  created_at: string
}

interface OmniScript {
  id: string
  category_id: string
  title: string
  content: string
  media_type: string
  file_name?: string
  delay_ms: number
  created_at: string
}

interface OmniSequenceStep {
  id: string
  sequence_id: string
  step_order: number
  title: string
  content: string
  media_type: string
  file_name?: string
  delay_ms: number
}

interface OmniSequence {
  id: string
  category_id: string
  title: string
  description: string
  steps: OmniSequenceStep[]
  created_at: string
}

const categories = ref<OmniCategory[]>([])
const sequences = ref<OmniSequence[]>([])
const isLoading = ref(true)
const activeTab = ref('categories')
const omnisSubTab = ref('scripts')

// Category dialog
const categoryDialogOpen = ref(false)
const isSubmittingCategory = ref(false)
const editingCategory = ref<OmniCategory | null>(null)
const categoryFormData = ref({
  name: '',
  color: '#3b82f6'
})

// Script dialog
const scriptDialogOpen = ref(false)
const isSubmittingScript = ref(false)
const editingScript = ref<OmniScript | null>(null)
const scriptFormData = ref({
  category_id: '',
  title: '',
  content: '',
  media_type: 'text',
  delay_ms: 0
})

// Sequence dialog
const sequenceDialogOpen = ref(false)
const isSubmittingSequence = ref(false)
const editingSequence = ref<OmniSequence | null>(null)
const sequenceFormData = ref({
  category_id: '',
  title: '',
  description: '',
  steps: [] as Array<{ title: string; content: string; media_type: string; delay_ms: number }>
})

// Delete dialogs
const deleteCategoryDialogOpen = ref(false)
const deleteScriptDialogOpen = ref(false)
const categoryToDelete = ref<OmniCategory | null>(null)
const scriptToDelete = ref<OmniScript | null>(null)

// Media upload
const mediaFile = ref<File | null>(null)
const isUploadingMedia = ref(false)

const mediaTypeOptions = [
  { value: 'text', label: 'Texto', icon: MessageSquare },
  { value: 'image', label: 'Imagem', icon: Image },
  { value: 'video', label: 'Vídeo', icon: Video },
  { value: 'audio', label: 'Áudio', icon: Mic },
  { value: 'document', label: 'Documento', icon: FileIcon }
]

const fetchCategories = async () => {
  isLoading.value = true
  try {
    const { data } = await api.get('/omnis/categories')
    categories.value = data || []
  } catch (error) {
    toast.error('Erro ao carregar categorias')
  } finally {
    isLoading.value = false
  }
}

const fetchSequences = async () => {
  try {
    const { data } = await api.get('/omnis/sequences')
    sequences.value = data || []
  } catch (error) {
    toast.error('Erro ao carregar sequências')
  }
}

// Category operations
const openCategoryDialog = (category?: OmniCategory) => {
  if (category) {
    editingCategory.value = category
    categoryFormData.value = {
      name: category.name,
      color: category.color
    }
  } else {
    editingCategory.value = null
    categoryFormData.value = {
      name: '',
      color: '#3b82f6'
    }
  }
  categoryDialogOpen.value = true
}

const saveCategory = async () => {
  if (!categoryFormData.value.name) {
    toast.error('Nome da categoria é obrigatório')
    return
  }

  isSubmittingCategory.value = true
  try {
    if (editingCategory.value) {
      await api.put(`/omnis/categories/${editingCategory.value.id}`, categoryFormData.value)
      toast.success('Categoria atualizada')
    } else {
      await api.post('/omnis/categories', categoryFormData.value)
      toast.success('Categoria criada')
    }
    categoryDialogOpen.value = false
    fetchCategories()
  } catch (error) {
    toast.error('Erro ao salvar categoria')
  } finally {
    isSubmittingCategory.value = false
  }
}

const confirmDeleteCategory = (category: OmniCategory) => {
  categoryToDelete.value = category
  deleteCategoryDialogOpen.value = true
}

const deleteCategory = async () => {
  if (!categoryToDelete.value) return

  try {
    await api.delete(`/omnis/categories/${categoryToDelete.value.id}`)
    toast.success('Categoria excluída')
    deleteCategoryDialogOpen.value = false
    fetchCategories()
  } catch (error) {
    toast.error('Erro ao excluir categoria')
  }
}

// Script operations
const openScriptDialog = (categoryId: string, script?: OmniScript) => {
  if (script) {
    editingScript.value = script
    scriptFormData.value = {
      category_id: script.category_id,
      title: script.title,
      content: script.content,
      media_type: script.media_type,
      delay_ms: script.delay_ms
    }
  } else {
    editingScript.value = null
    scriptFormData.value = {
      category_id: categoryId,
      title: '',
      content: '',
      media_type: 'text',
      delay_ms: 0
    }
  }
  mediaFile.value = null
  scriptDialogOpen.value = true
}

const saveScript = async () => {
  if (!scriptFormData.value.title) {
    toast.error('Título da mensagem é obrigatório')
    return
  }

  isSubmittingScript.value = true
  try {
    let scriptId: string

    if (editingScript.value) {
      await api.put(`/omnis/scripts/${editingScript.value.id}`, scriptFormData.value)
      scriptId = editingScript.value.id
      toast.success('Mensagem atualizada')
    } else {
      const { data } = await api.post('/omnis/scripts', scriptFormData.value)
      scriptId = data.id
      toast.success('Mensagem criada')
    }

    // Upload media if file is selected
    if (mediaFile.value && scriptFormData.value.media_type !== 'text') {
      await uploadMedia(scriptId)
    }

    scriptDialogOpen.value = false
    fetchCategories()
  } catch (error) {
    toast.error('Erro ao salvar mensagem')
  } finally {
    isSubmittingScript.value = false
  }
}

const uploadMedia = async (scriptId: string) => {
  if (!mediaFile.value) return

  isUploadingMedia.value = true
  try {
    const formData = new FormData()
    formData.append('file', mediaFile.value)

    await api.post(`/omnis/scripts/${scriptId}/media`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    toast.success('Mídia enviada')
  } catch (error) {
    toast.error('Erro ao enviar mídia')
  } finally {
    isUploadingMedia.value = false
  }
}

const confirmDeleteScript = (script: OmniScript) => {
  scriptToDelete.value = script
  deleteScriptDialogOpen.value = true
}

const deleteScript = async () => {
  if (!scriptToDelete.value) return

  try {
    await api.delete(`/omnis/scripts/${scriptToDelete.value.id}`)
    toast.success('Mensagem excluída')
    deleteScriptDialogOpen.value = false
    fetchCategories()
  } catch (error) {
    toast.error('Erro ao excluir mensagem')
  }
}

const handleMediaSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    mediaFile.value = target.files[0]
  }
}

const getMediaTypeIcon = (type: string) => {
  const option = mediaTypeOptions.find(o => o.value === type)
  return option?.icon || MessageSquare
}

const getMediaTypeLabel = (type: string) => {
  const option = mediaTypeOptions.find(o => o.value === type)
  return option?.label || 'Texto'
}

// Export/Import operations
const isExporting = ref(false)
const isImporting = ref(false)
const importFile = ref<File | null>(null)
const importDialogOpen = ref(false)

const exportOmnisJSON = async () => {
  isExporting.value = true
  try {
    const { data } = await api.post('/omnis/export', { include_media: true })
    
    // Create and download JSON file
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `omnis-export-${new Date().toISOString().split('T')[0]}.json`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    
    toast.success('Exportação concluída')
  } catch (error) {
    toast.error('Erro ao exportar')
  } finally {
    isExporting.value = false
  }
}

const exportOmnisZIP = async () => {
  isExporting.value = true
  try {
    const response = await api.get('/omnis/export/zip', {
      responseType: 'blob'
    })
    
    // Create and download ZIP file
    const blob = new Blob([response.data], { type: 'application/zip' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `omnis-export-${new Date().toISOString().split('T')[0]}.zip`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    
    toast.success('Exportação ZIP concluída')
  } catch (error) {
    toast.error('Erro ao exportar ZIP')
  } finally {
    isExporting.value = false
  }
}

const openImportDialog = () => {
  importFile.value = null
  importDialogOpen.value = true
}

const handleImportFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    importFile.value = target.files[0]
  }
}

const importOmnisJSON = async () => {
  if (!importFile.value) {
    toast.error('Selecione um arquivo')
    return
  }

  isImporting.value = true
  try {
    const text = await importFile.value.text()
    const data = JSON.parse(text)
    
    await api.post('/omnis/import', {
      replace_existing: true,
      data: data
    })
    
    toast.success('Importação concluída')
    importDialogOpen.value = false
    fetchCategories()
    fetchSequences()
  } catch (error) {
    toast.error('Erro ao importar JSON')
  } finally {
    isImporting.value = false
  }
}

const importOmnisZIP = async () => {
  if (!importFile.value) {
    toast.error('Selecione um arquivo')
    return
  }

  isImporting.value = true
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)
    
    await api.post('/omnis/import/zip', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    toast.success('Importação ZIP concluída')
    importDialogOpen.value = false
    fetchCategories()
    fetchSequences()
  } catch (error) {
    toast.error('Erro ao importar ZIP')
  } finally {
    isImporting.value = false
  }
}

// Sequence operations
const openSequenceDialog = (categoryId?: string, sequence?: OmniSequence) => {
  if (sequence) {
    editingSequence.value = sequence
    sequenceFormData.value = {
      category_id: sequence.category_id,
      title: sequence.title,
      description: sequence.description || '',
      steps: sequence.steps?.map(s => ({
        title: s.title,
        content: s.content,
        media_type: s.media_type,
        delay_ms: s.delay_ms
      })) || []
    }
  } else {
    editingSequence.value = null
    sequenceFormData.value = {
      category_id: categoryId || '',
      title: '',
      description: '',
      steps: []
    }
  }
  sequenceDialogOpen.value = true
}

const addSequenceStep = () => {
  sequenceFormData.value.steps.push({
    title: '',
    content: '',
    media_type: 'text',
    delay_ms: 1000
  })
}

const removeSequenceStep = (index: number) => {
  sequenceFormData.value.steps.splice(index, 1)
}

const moveStepUp = (index: number) => {
  if (index > 0) {
    const steps = sequenceFormData.value.steps
    const temp = steps[index]
    steps[index] = steps[index - 1]
    steps[index - 1] = temp
  }
}

const moveStepDown = (index: number) => {
  const steps = sequenceFormData.value.steps
  if (index < steps.length - 1) {
    const temp = steps[index]
    steps[index] = steps[index + 1]
    steps[index + 1] = temp
  }
}

const saveSequence = async () => {
  if (!sequenceFormData.value.title) {
    toast.error('Título da sequência é obrigatório')
    return
  }

  if (sequenceFormData.value.steps.length === 0) {
    toast.error('Adicione pelo menos uma mensagem à sequência')
    return
  }

  isSubmittingSequence.value = true
  try {
    if (editingSequence.value) {
      await api.put(`/omnis/sequences/${editingSequence.value.id}`, {
        title: sequenceFormData.value.title,
        description: sequenceFormData.value.description
      })
      toast.success('Sequência atualizada')
    } else {
      await api.post('/omnis/sequences', sequenceFormData.value)
      toast.success('Sequência criada')
    }
    sequenceDialogOpen.value = false
    fetchSequences()
  } catch (error) {
    toast.error('Erro ao salvar sequência')
  } finally {
    isSubmittingSequence.value = false
  }
}

const sendSequence = async (sequenceId: string, contactId: string) => {
  try {
    await api.post('/omnis/send-sequence', {
      sequence_id: sequenceId,
      contact_id: contactId
    })
    toast.success('Sequência enviada')
  } catch (error) {
    toast.error('Erro ao enviar sequência')
  }
}

onMounted(() => {
  fetchCategories()
  fetchSequences()
})
</script>

<template>
  <div class="p-6 space-y-6">
    <PageHeader
      title="Mensagens Rápidas (Omni)"
      description="Crie e gerencie mensagens pré-definidas para envio rápido"
    >
      <div class="flex items-center gap-2">
        <Button variant="outline" @click="openImportDialog" class="gap-2">
          <Import class="w-4 h-4" />
          Importar
        </Button>
        <Button variant="outline" @click="exportOmnisJSON" :disabled="isExporting" class="gap-2">
          <Download class="w-4 h-4" />
          Exportar JSON
        </Button>
        <Button variant="outline" @click="exportOmnisZIP" :disabled="isExporting" class="gap-2">
          <Download class="w-4 h-4" />
          Exportar ZIP
        </Button>
        <Button @click="openCategoryDialog()" class="gap-2">
          <Plus class="w-4 h-4" />
          Nova Categoria
        </Button>
      </div>
    </PageHeader>

    <Tabs v-model="activeTab" class="space-y-4">
      <TabsList>
        <TabsTrigger value="categories">Categorias e Mensagens</TabsTrigger>
        <TabsTrigger value="sequences">Fluxos de Mensagens</TabsTrigger>
        <TabsTrigger value="import">Importar ZIP</TabsTrigger>
      </TabsList>

      <TabsContent value="categories" class="space-y-4">
        <div v-if="isLoading" class="flex items-center justify-center py-12">
          <Loader2 class="w-8 h-8 animate-spin text-muted-foreground" />
        </div>

        <div v-else-if="categories.length === 0" class="text-center py-12">
          <MessageSquare class="w-12 h-12 mx-auto text-muted-foreground mb-4" />
          <h3 class="text-lg font-medium">Nenhuma categoria criada</h3>
          <p class="text-muted-foreground mt-1">Crie categorias para organizar suas mensagens rápidas</p>
          <Button @click="openCategoryDialog()" class="mt-4 gap-2">
            <Plus class="w-4 h-4" />
            Criar Primeira Categoria
          </Button>
        </div>

        <div v-else class="grid gap-4">
          <Card v-for="category in categories" :key="category.id">
            <CardHeader class="pb-3">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div
                    class="w-4 h-4 rounded-full"
                    :style="{ backgroundColor: category.color }"
                  />
                  <CardTitle>{{ category.name }}</CardTitle>
                  <Badge variant="secondary">{{ category.scripts?.length || 0 }} mensagens</Badge>
                </div>
                <div class="flex items-center gap-2">
                  <Button variant="ghost" size="icon" @click="openScriptDialog(category.id)">
                    <Plus class="w-4 h-4" />
                  </Button>
                  <Button variant="ghost" size="icon" @click="openCategoryDialog(category)">
                    <Pencil class="w-4 h-4" />
                  </Button>
                  <Button variant="ghost" size="icon" @click="confirmDeleteCategory(category)">
                    <Trash2 class="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div v-if="!category.scripts?.length" class="text-center py-6 text-muted-foreground">
                <p>Nenhuma mensagem nesta categoria</p>
                <Button variant="outline" size="sm" @click="openScriptDialog(category.id)" class="mt-2">
                  <Plus class="w-4 h-4 mr-1" />
                  Adicionar Mensagem
                </Button>
              </div>

              <ScrollArea v-else class="max-h-64">
                <div class="space-y-2">
                  <div
                    v-for="script in category.scripts"
                    :key="script.id"
                    class="flex items-center justify-between p-3 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
                  >
                    <div class="flex items-center gap-3 min-w-0">
                      <component
                        :is="getMediaTypeIcon(script.media_type)"
                        class="w-4 h-4 text-muted-foreground flex-shrink-0"
                      />
                      <div class="min-w-0">
                        <p class="font-medium truncate">{{ script.title }}</p>
                        <p class="text-sm text-muted-foreground truncate">
                          {{ script.media_type === 'text' ? script.content : script.file_name || 'Arquivo de mídia' }}
                        </p>
                      </div>
                    </div>
                    <div class="flex items-center gap-1">
                      <Badge variant="outline" class="text-xs">
                        {{ getMediaTypeLabel(script.media_type) }}
                      </Badge>
                      <Button variant="ghost" size="icon" @click="openScriptDialog(category.id, script)">
                        <Pencil class="w-4 h-4" />
                      </Button>
                      <Button variant="ghost" size="icon" @click="confirmDeleteScript(script)">
                        <Trash2 class="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              </ScrollArea>
            </CardContent>
          </Card>
        </div>
      </TabsContent>

      <TabsContent value="sequences" class="space-y-4">
        <div v-if="categories.length === 0" class="text-center py-12">
          <MessageSquare class="w-12 h-12 mx-auto text-muted-foreground mb-4" />
          <h3 class="text-lg font-medium">Crie categorias primeiro</h3>
          <p class="text-muted-foreground mt-1">Você precisa ter pelo menos uma categoria para criar fluxos de mensagens</p>
        </div>

        <div v-else class="flex justify-between items-center mb-4">
          <p class="text-muted-foreground">
            Crie sequências de mensagens para enviar múltiplos conteúdos com 1 clique
          </p>
          <Button @click="openSequenceDialog()" class="gap-2">
            <Plus class="w-4 h-4" />
            Novo Fluxo
          </Button>
        </div>

        <div v-if="sequences.length === 0 && categories.length > 0" class="text-center py-12 border-2 border-dashed rounded-lg">
          <MessageSquare class="w-12 h-12 mx-auto text-muted-foreground mb-4" />
          <h3 class="text-lg font-medium">Nenhum fluxo criado</h3>
          <p class="text-muted-foreground mt-1">Crie fluxos de mensagens para envio sequencial</p>
          <Button @click="openSequenceDialog()" class="mt-4 gap-2">
            <Plus class="w-4 h-4" />
            Criar Primeiro Fluxo
          </Button>
        </div>

        <div v-else class="grid gap-4">
          <Card v-for="sequence in sequences" :key="sequence.id">
            <CardHeader class="pb-3">
              <div class="flex items-center justify-between">
                <div>
                  <CardTitle>{{ sequence.title }}</CardTitle>
                  <CardDescription v-if="sequence.description">{{ sequence.description }}</CardDescription>
                </div>
                <div class="flex items-center gap-2">
                  <Badge variant="secondary">{{ sequence.steps?.length || 0 }} mensagens</Badge>
                  <Button variant="outline" size="sm" @click="openSequenceDialog(undefined, sequence)">
                    <Pencil class="w-4 h-4 mr-1" />
                    Editar
                  </Button>
                  <Button variant="destructive" size="sm" @click="async () => {
                    await api.delete(`/omnis/sequences/${sequence.id}`)
                    toast.success('Fluxo excluído')
                    fetchSequences()
                  }">
                    <Trash2 class="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div v-if="!sequence.steps?.length" class="text-center py-4 text-muted-foreground">
                <p>Nenhum passo adicionado</p>
              </div>

              <ScrollArea v-else class="max-h-48">
                <div class="space-y-2">
                  <div
                    v-for="(step, index) in sequence.steps"
                    :key="step.id"
                    class="flex items-center gap-3 p-3 rounded-lg border bg-card"
                  >
                    <div class="flex items-center justify-center w-6 h-6 rounded-full bg-primary text-primary-foreground text-sm font-medium">
                      {{ index + 1 }}
                    </div>
                    <component
                      :is="getMediaTypeIcon(step.media_type)"
                      class="w-4 h-4 text-muted-foreground"
                    />
                    <div class="flex-1 min-w-0">
                      <p class="font-medium truncate">{{ step.title }}</p>
                      <p class="text-sm text-muted-foreground truncate">
                        {{ step.media_type === 'text' ? step.content : step.file_name || 'Mídia' }}
                      </p>
                    </div>
                    <Badge variant="outline" class="text-xs">
                      {{ step.delay_ms }}ms
                    </Badge>
                  </div>
                </div>
              </ScrollArea>
            </CardContent>
          </Card>
        </div>
      </TabsContent>

      <TabsContent value="import" class="space-y-4">
        <Card>
          <CardHeader>
            <CardTitle>Importar via ZIP</CardTitle>
            <CardDescription>
              Importe um arquivo ZIP contendo config.json e arquivos de mídia
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div class="border-2 border-dashed rounded-lg p-8 text-center">
              <Upload class="w-12 h-12 mx-auto text-muted-foreground mb-4" />
              <p class="text-muted-foreground mb-4">
                Arraste um arquivo ZIP ou clique para selecionar
              </p>
              <input
                type="file"
                accept=".zip"
                class="hidden"
                id="zip-upload"
                @change="async (e) => {
                  const file = (e.target as HTMLInputElement).files?.[0]
                  if (!file) return

                  const formData = new FormData()
                  formData.append('file', file)

                  try {
                    await api.post('/omnis/import', formData, {
                      headers: { 'Content-Type': 'multipart/form-data' }
                    })
                    toast.success('Importação concluída')
                    fetchCategories()
                  } catch (error) {
                    toast.error('Erro na importação')
                  }
                }"
              />
              <Button variant="outline" @click="() => document.getElementById('zip-upload')?.click()">
                <Upload class="w-4 h-4 mr-2" />
                Selecionar Arquivo
              </Button>
            </div>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>

    <!-- Category Dialog -->
    <Dialog v-model:open="categoryDialogOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ editingCategory ? 'Editar Categoria' : 'Nova Categoria' }}</DialogTitle>
          <DialogDescription>
            {{ editingCategory ? 'Atualize as informações da categoria' : 'Crie uma nova categoria para organizar mensagens' }}
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>Nome da Categoria</Label>
            <Input v-model="categoryFormData.name" placeholder="Ex: Vendas, Suporte, etc." />
          </div>
          <div class="space-y-2">
            <Label>Cor</Label>
            <div class="flex items-center gap-3">
              <input
                type="color"
                v-model="categoryFormData.color"
                class="w-12 h-10 rounded cursor-pointer"
              />
              <Input v-model="categoryFormData.color" placeholder="#3b82f6" />
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="categoryDialogOpen = false">Cancelar</Button>
          <Button @click="saveCategory" :disabled="isSubmittingCategory">
            <Loader2 v-if="isSubmittingCategory" class="w-4 h-4 mr-2 animate-spin" />
            {{ editingCategory ? 'Salvar' : 'Criar' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Script Dialog -->
    <Dialog v-model:open="scriptDialogOpen">
      <DialogContent class="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{{ editingScript ? 'Editar Mensagem' : 'Nova Mensagem' }}</DialogTitle>
          <DialogDescription>
            {{ editingScript ? 'Atualize a mensagem pré-definida' : 'Crie uma nova mensagem para envio rápido' }}
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>Título da Mensagem</Label>
            <Input v-model="scriptFormData.title" placeholder="Ex: Mensagem de boas-vindas" />
          </div>

          <div class="space-y-2">
            <Label>Tipo de Mídia</Label>
            <Select v-model="scriptFormData.media_type">
              <SelectTrigger>
                <SelectValue placeholder="Selecione o tipo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="option in mediaTypeOptions" :key="option.value" :value="option.value">
                  <div class="flex items-center gap-2">
                    <component :is="option.icon" class="w-4 h-4" />
                    {{ option.label }}
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div v-if="scriptFormData.media_type === 'text'" class="space-y-2">
            <Label>Conteúdo da Mensagem</Label>
            <Textarea
              v-model="scriptFormData.content"
              placeholder="Digite a mensagem..."
              rows="4"
            />
          </div>

          <div v-else class="space-y-2">
            <Label>{{ scriptFormData.media_type === 'audio' ? 'Arquivo de Áudio' : 'Arquivo de Mídia' }}</Label>
            <div class="border-2 border-dashed rounded-lg p-4 text-center">
              <input
                type="file"
                :accept="scriptFormData.media_type === 'audio' ? 'audio/*' : scriptFormData.media_type + '/*'"
                class="hidden"
                id="media-upload"
                @change="handleMediaSelect"
              />
              <div v-if="mediaFile" class="flex items-center justify-center gap-2">
                <component :is="getMediaTypeIcon(scriptFormData.media_type)" class="w-4 h-4" />
                <span class="text-sm">{{ mediaFile.name }}</span>
                <Button variant="ghost" size="icon" @click="mediaFile = null">
                  <X class="w-4 h-4" />
                </Button>
              </div>
              <div v-else>
                <p class="text-sm text-muted-foreground mb-2">Selecione um arquivo</p>
                <Button variant="outline" size="sm" @click="() => document.getElementById('media-upload')?.click()">
                  <Upload class="w-4 h-4 mr-1" />
                  Selecionar
                </Button>
              </div>
            </div>
            <div class="space-y-2">
              <Label>Legenda (opcional)</Label>
              <Textarea
                v-model="scriptFormData.content"
                placeholder="Legenda para a mídia..."
                rows="2"
              />
            </div>
          </div>

          <div class="space-y-2">
            <Label>Delay antes do envio (ms)</Label>
            <Input
              v-model.number="scriptFormData.delay_ms"
              type="number"
              min="0"
              placeholder="0"
            />
            <p class="text-sm text-muted-foreground">
              Tempo de espera antes de enviar esta mensagem
            </p>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="scriptDialogOpen = false">Cancelar</Button>
          <Button @click="saveScript" :disabled="isSubmittingScript">
            <Loader2 v-if="isSubmittingScript" class="w-4 h-4 mr-2 animate-spin" />
            {{ editingScript ? 'Salvar' : 'Criar' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Category Dialog -->
    <AlertDialog v-model:open="deleteCategoryDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Excluir Categoria</AlertDialogTitle>
          <AlertDialogDescription>
            Tem certeza que deseja excluir a categoria "{{ categoryToDelete?.name }}"?
            Todas as mensagens nesta categoria serão excluídas.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancelar</AlertDialogCancel>
          <AlertDialogAction @click="deleteCategory">Excluir</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- Delete Script Dialog -->
    <AlertDialog v-model:open="deleteScriptDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Excluir Mensagem</AlertDialogTitle>
          <AlertDialogDescription>
            Tem certeza que deseja excluir a mensagem "{{ scriptToDelete?.title }}"?
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancelar</AlertDialogCancel>
          <AlertDialogAction @click="deleteScript">Excluir</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- Sequence Dialog -->
    <Dialog v-model:open="sequenceDialogOpen">
      <DialogContent class="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{{ editingSequence ? 'Editar Fluxo' : 'Novo Fluxo de Mensagens' }}</DialogTitle>
          <DialogDescription>
            {{ editingSequence ? 'Edite a sequência de mensagens' : 'Crie uma sequência de mensagens para envio rápido' }}
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-6 py-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label>Categoria</Label>
              <Select v-model="sequenceFormData.category_id" :disabled="!!editingSequence">
                <SelectTrigger>
                  <SelectValue placeholder="Selecione a categoria" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="cat in categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="space-y-2">
              <Label>Título do Fluxo</Label>
              <Input v-model="sequenceFormData.title" placeholder="Ex: Processo de Vendas" />
            </div>
          </div>

          <div class="space-y-2">
            <Label>Descrição (opcional)</Label>
            <Textarea
              v-model="sequenceFormData.description"
              placeholder="Descreva o propósito deste fluxo..."
              rows="2"
            />
          </div>

          <Separator />

          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <Label class="text-lg">Mensagens do Fluxo</Label>
              <Button variant="outline" size="sm" @click="addSequenceStep" class="gap-1">
                <Plus class="w-4 h-4" />
                Adicionar Mensagem
              </Button>
            </div>

            <div v-if="sequenceFormData.steps.length === 0" class="text-center py-8 border-2 border-dashed rounded-lg">
              <MessageSquare class="w-8 h-8 mx-auto text-muted-foreground mb-2" />
              <p class="text-muted-foreground">Adicione mensagens ao fluxo</p>
            </div>

            <div v-else class="space-y-4">
              <Card v-for="(step, index) in sequenceFormData.steps" :key="index" class="relative">
                <CardHeader class="pb-2">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-2">
                      <div class="flex items-center justify-center w-6 h-6 rounded-full bg-primary text-primary-foreground text-sm font-medium">
                        {{ index + 1 }}
                      </div>
                      <Input
                        v-model="step.title"
                        placeholder="Título da mensagem"
                        class="w-48"
                      />
                    </div>
                    <div class="flex items-center gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        @click="moveStepUp(index)"
                        :disabled="index === 0"
                      >
                        ↑
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        @click="moveStepDown(index)"
                        :disabled="index === sequenceFormData.steps.length - 1"
                      >
                        ↓
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        @click="removeSequenceStep(index)"
                      >
                        <X class="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </CardHeader>
                <CardContent class="space-y-3">
                  <div class="grid grid-cols-2 gap-3">
                    <div class="space-y-2">
                      <Label class="text-sm">Tipo</Label>
                      <Select v-model="step.media_type">
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem v-for="option in mediaTypeOptions" :key="option.value" :value="option.value">
                            <div class="flex items-center gap-2">
                              <component :is="option.icon" class="w-4 h-4" />
                              {{ option.label }}
                            </div>
                          </SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div class="space-y-2">
                      <Label class="text-sm">Delay (ms)</Label>
                      <Input v-model.number="step.delay_ms" type="number" min="0" />
                    </div>
                  </div>

                  <div v-if="step.media_type === 'text'" class="space-y-2">
                    <Label class="text-sm">Conteúdo</Label>
                    <Textarea v-model="step.content" placeholder="Digite a mensagem..." rows="2" />
                  </div>

                  <div v-else class="space-y-2">
                    <Label class="text-sm">Arquivo / Legenda</Label>
                    <Input v-model="step.content" placeholder="Caminho do arquivo ou legenda" />
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="sequenceDialogOpen = false">Cancelar</Button>
          <Button @click="saveSequence" :disabled="isSubmittingSequence">
            <Loader2 v-if="isSubmittingSequence" class="w-4 h-4 mr-2 animate-spin" />
            {{ editingSequence ? 'Salvar' : 'Criar Fluxo' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Import Dialog -->
    <Dialog v-model:open="importDialogOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Importar Mensagens Rápidas</DialogTitle>
          <DialogDescription>
            Importe um arquivo JSON ou ZIP com categorias, mensagens e fluxos
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="border-2 border-dashed rounded-lg p-8 text-center">
            <Import class="w-12 h-12 mx-auto text-muted-foreground mb-4" />
            <p class="text-muted-foreground mb-4">
              Selecione um arquivo JSON ou ZIP exportado anteriormente
            </p>
            <input
              type="file"
              accept=".json,.zip"
              class="hidden"
              id="omnis-import-file"
              @change="handleImportFileSelect"
            />
            <Button variant="outline" @click="() => document.getElementById('omnis-import-file')?.click()">
              <Upload class="w-4 h-4 mr-2" />
              Selecionar Arquivo
            </Button>
          </div>

          <div v-if="importFile" class="flex items-center gap-2 p-3 bg-muted rounded-lg">
            <FileIcon class="w-4 h-4 text-muted-foreground" />
            <span class="text-sm flex-1">{{ importFile.name }}</span>
            <span class="text-xs text-muted-foreground">
              {{ (importFile.size / 1024).toFixed(1) }} KB
            </span>
          </div>

          <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-3">
            <p class="text-sm text-yellow-800 dark:text-yellow-200">
              <strong>Atenção:</strong> A importação irá substituir todas as mensagens rápidas existentes.
            </p>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="importDialogOpen = false">Cancelar</Button>
          <Button 
            @click="importFile?.name.endsWith('.zip') ? importOmnisZIP() : importOmnisJSON()" 
            :disabled="isImporting || !importFile"
          >
            <Loader2 v-if="isImporting" class="w-4 h-4 mr-2 animate-spin" />
            Importar
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
