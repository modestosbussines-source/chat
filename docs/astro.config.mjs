import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: 'https://omni.com.br',
  base: '/docs',
  integrations: [
    starlight({
      title: 'OMNI',
      description: 'Documentação da Plataforma OMNI - Atendimento Omnichannel',
      social: {},
      customCss: [
        './src/styles/custom.css',
      ],
      sidebar: [
        {
          label: 'Introdução',
          items: [
            { label: 'Sobre o OMNI', slug: 'getting-started/introduction' },
            { label: 'Primeiros Passos', slug: 'getting-started/quickstart' },
            { label: 'Configuração', slug: 'getting-started/configuration' },
          ],
        },
        {
          label: 'Recursos',
          items: [
            { label: 'Dashboard', slug: 'features/dashboard' },
            { label: 'Funções e Permissões', slug: 'features/roles-permissions' },
            { label: 'Chatbot com IA', slug: 'features/chatbot' },
            { label: 'Respostas Rápidas', slug: 'features/canned-responses' },
            { label: 'Ações Personalizadas', slug: 'features/custom-actions' },
            { label: 'Templates de Mensagem', slug: 'features/templates' },
            { label: 'Campanhas', slug: 'features/campaigns' },
            { label: 'WhatsApp Flows', slug: 'features/whatsapp-flows' },
            { label: 'Chamadas', slug: 'features/calling' },
          ],
        },
        {
          label: 'Integrações',
          items: [
            { label: 'Visão Geral da API', slug: 'api-reference/overview' },
            { label: 'Autenticação', slug: 'api-reference/authentication' },
            { label: 'Chaves de API', slug: 'api-reference/api-keys' },
            { label: 'Usuários', slug: 'api-reference/users' },
            { label: 'Organizações', slug: 'api-reference/organizations' },
            { label: 'Funções', slug: 'api-reference/roles' },
            { label: 'Contas', slug: 'api-reference/accounts' },
            { label: 'Contatos', slug: 'api-reference/contacts' },
            { label: 'Mensagens', slug: 'api-reference/messages' },
            { label: 'Templates', slug: 'api-reference/templates' },
            { label: 'Flows', slug: 'api-reference/flows' },
            { label: 'Campanhas', slug: 'api-reference/campaigns' },
            { label: 'Chatbot', slug: 'api-reference/chatbot' },
            { label: 'Respostas Rápidas', slug: 'api-reference/canned-responses' },
            { label: 'Ações Customizadas', slug: 'api-reference/custom-actions' },
            { label: 'Webhooks', slug: 'api-reference/webhooks' },
            { label: 'Analytics', slug: 'api-reference/analytics' },
          ],
        },
      ],
    }),
  ],
});