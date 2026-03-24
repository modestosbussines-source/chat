/** @type {import('tailwindcss').Config} */
const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  darkMode: ["class"],
  plugins: [
    require("tailwindcss-animate"),
    function({ addVariant }) {
      addVariant('light', '.light &')
    }
  ],
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px'
      }
    },
    extend: {
      fontFamily: {
        sans: ['Inter', ...defaultTheme.fontFamily.sans],
        mono: ['JetBrains Mono', ...defaultTheme.fontFamily.mono],
      },
      colors: {
        // Core semantic colors
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        
        // Primary - Indigo/Violet
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
          50: 'hsl(252, 100%, 98%)',
          100: 'hsl(252, 100%, 95%)',
          200: 'hsl(252, 95%, 90%)',
          300: 'hsl(252, 85%, 78%)',
          400: 'hsl(252, 75%, 65%)',
          500: 'hsl(252, 70%, 55%)',
          600: 'hsl(252, 70%, 48%)',
          700: 'hsl(252, 68%, 40%)',
          800: 'hsl(252, 65%, 32%)',
          900: 'hsl(252, 60%, 24%)',
        },
        
        // Secondary - Teal/Cyan
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
          50: 'hsl(199, 100%, 97%)',
          100: 'hsl(199, 100%, 94%)',
          200: 'hsl(199, 90%, 88%)',
          300: 'hsl(199, 80%, 76%)',
          400: 'hsl(199, 75%, 62%)',
          500: 'hsl(199, 80%, 50%)',
          600: 'hsl(199, 80%, 42%)',
          700: 'hsl(199, 78%, 34%)',
          800: 'hsl(199, 75%, 26%)',
          900: 'hsl(199, 70%, 18%)',
        },
        
        // Semantic colors
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        success: {
          DEFAULT: 'hsl(var(--success))',
          foreground: 'hsl(var(--success-foreground))',
        },
        warning: {
          DEFAULT: 'hsl(var(--warning))',
          foreground: 'hsl(var(--warning-foreground))',
        },
        info: {
          DEFAULT: 'hsl(var(--info))',
          foreground: 'hsl(var(--info-foreground))',
        },
        
        // Muted/Accent
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        
        // Overlays
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        
        // Accent colors - Brazilian vibrancy
        coral: {
          light: 'hsl(12, 100%, 72%)',
          DEFAULT: 'hsl(12, 100%, 62%)',
          dark: 'hsl(12, 100%, 52%)',
        },
        teal: {
          light: 'hsl(168, 80%, 65%)',
          DEFAULT: 'hsl(168, 80%, 50%)',
          dark: 'hsl(168, 80%, 38%)',
        },
        amber: {
          light: 'hsl(38, 92%, 70%)',
          DEFAULT: 'hsl(38, 92%, 58%)',
          dark: 'hsl(38, 92%, 46%)',
        },
        pink: {
          light: 'hsl(330, 100%, 78%)',
          DEFAULT: 'hsl(330, 100%, 68%)',
          dark: 'hsl(330, 100%, 58%)',
        },
        
        // WhatsApp brand
        whatsapp: {
          green: '#25D366',
          teal: '#128C7E',
          'teal-dark': '#075E54',
          light: '#DCF8C6'
        },
        
        // Glass effect colors
        glass: {
          bg: 'var(--glass-bg)',
          'bg-hover': 'var(--glass-bg-hover)',
          border: 'var(--glass-border)',
        },
        
        // Chat-specific colors
        chat: {
          bg: 'hsl(var(--chat-bg))',
          'bubble-incoming': 'hsl(var(--chat-bubble-incoming))',
          'bubble-outgoing': 'hsl(var(--chat-bubble-outgoing))',
        },
        
        // Inbox-specific colors
        inbox: {
          sidebar: 'hsl(var(--inbox-sidebar))',
          'item-hover': 'hsl(var(--inbox-item-hover))',
          'item-active': 'hsl(var(--inbox-item-active))',
          'unread-badge': 'hsl(var(--inbox-unread-badge))',
        },
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
        xl: 'var(--radius, 0.75rem)',
        '2xl': 'calc(var(--radius, 0.75rem) + 4px)',
      },
      boxShadow: {
        'sm': 'var(--shadow-sm)',
        'md': 'var(--shadow-md)',
        'lg': 'var(--shadow-lg)',
        'xl': 'var(--shadow-xl)',
      },
      keyframes: {
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--reka-accordion-content-height)' }
        },
        'accordion-up': {
          from: { height: 'var(--reka-accordion-content-height)' },
          to: { height: '0' }
        },
        'fade-in': {
          from: { opacity: '0', transform: 'translateY(4px)' },
          to: { opacity: '1', transform: 'translateY(0)' }
        },
        'slide-in-right': {
          from: { opacity: '0', transform: 'translateX(10px)' },
          to: { opacity: '1', transform: 'translateX(0)' }
        },
        'pulse-ring': {
          '0%': { transform: 'scale(0.95)', opacity: '0.7' },
          '50%': { transform: 'scale(1.05)', opacity: '1' },
          '100%': { transform: 'scale(0.95)', opacity: '0.7' }
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
        'fade-in': 'fade-in 0.2s ease-out',
        'slide-in-right': 'slide-in-right 0.2s ease-out',
        'pulse-ring': 'pulse-ring 2s ease-in-out infinite',
      },
      backdropBlur: {
        xs: '2px',
      },
      spacing: {
        '18': '4.5rem',
        '72': '18rem',
        '80': '20rem',
      },
      transitionDuration: {
        '50': '50ms',
        '100': '100ms',
        '150': '150ms',
        '200': '200ms',
        '300': '300ms',
        '500': '500ms',
      },
      transitionTimingFunction: {
        'bounce': 'cubic-bezier(0.68, -0.55, 0.265, 1.55)',
      },
      fontSize: {
        '2xs': ['0.625rem', { lineHeight: '1rem' }],
      },
      zIndex: {
        '60': '60',
        '70': '70',
        '80': '80',
        '90': '90',
        '100': '100',
      },
    }
  },
}
