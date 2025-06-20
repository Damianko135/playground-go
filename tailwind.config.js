/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.{templ,go}",
    "./cmd/**/*.go",
    "./internal/**/*.go",
    "./**/*.html"
  ],
  theme: {
    extend: {
      colors: {
        // Custom green palette
        primary: {
          50: '#ecfdf5',
          100: '#d1fae5',
          200: '#a7f3d0',
          300: '#6ee7b7',
          400: '#34d399',
          500: '#10b981',
          600: '#059669',
          700: '#047857',
          800: '#065f46',
          900: '#064e3b',
          950: '#022c22'
        },
        // Extended green variants
        forest: {
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0',
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e',
          600: '#16a34a',
          700: '#15803d',
          800: '#166534',
          900: '#14532d',
          950: '#052e16'
        },
        mint: {
          50: '#f0fdfa',
          100: '#ccfbf1',
          200: '#99f6e4',
          300: '#5eead4',
          400: '#2dd4bf',
          500: '#14b8a6',
          600: '#0d9488',
          700: '#0f766e',
          800: '#115e59',
          900: '#134e4a',
          950: '#042f2e'
        },
        sage: {
          50: '#f6f7f6',
          100: '#e3e8e3',
          200: '#c7d2c7',
          300: '#a3b5a3',
          400: '#7a9379',
          500: '#5a7659',
          600: '#465d45',
          700: '#3a4b39',
          800: '#313e30',
          900: '#2a342a',
          950: '#161b16'
        }
      },
      backgroundColor: {
        'app': '#f0fdf4',
        'card': '#ffffff',
        'card-hover': '#f9fffe'
      },
      textColor: {
        'app-primary': '#064e3b',
        'app-secondary': '#047857',
        'app-muted': '#6b7280'
      },
      borderColor: {
        'app-light': '#d1fae5',
        'app-medium': '#a7f3d0',
        'app-dark': '#059669'
      },
      ringColor: {
        'app': '#10b981'
      },
      gradientColorStops: {
        'green-start': '#ecfdf5',
        'green-middle': '#a7f3d0',
        'green-end': '#047857'
      },
      boxShadow: {
        'green-sm': '0 1px 2px 0 rgba(16, 185, 129, 0.05)',
        'green': '0 1px 3px 0 rgba(16, 185, 129, 0.1), 0 1px 2px 0 rgba(16, 185, 129, 0.06)',
        'green-md': '0 4px 6px -1px rgba(16, 185, 129, 0.1), 0 2px 4px -1px rgba(16, 185, 129, 0.06)',
        'green-lg': '0 10px 15px -3px rgba(16, 185, 129, 0.1), 0 4px 6px -2px rgba(16, 185, 129, 0.05)',
        'green-xl': '0 20px 25px -5px rgba(16, 185, 129, 0.1), 0 10px 10px -5px rgba(16, 185, 129, 0.04)'
      },
      animation: {
        'pulse-green': 'pulse-green 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'bounce-green': 'bounce-green 1s infinite'
      },
      keyframes: {
        'pulse-green': {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '.7', backgroundColor: '#a7f3d0' }
        },
        'bounce-green': {
          '0%, 100%': {
            transform: 'translateY(-25%)',
            animationTimingFunction: 'cubic-bezier(0.8, 0, 1, 1)'
          },
          '50%': {
            transform: 'translateY(0)',
            animationTimingFunction: 'cubic-bezier(0, 0, 0.2, 1)'
          }
        }
      },
      fontFamily: {
        'sans': ['Inter', 'system-ui', 'sans-serif'],
        'display': ['Poppins', 'system-ui', 'sans-serif']
      },
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
        '128': '32rem'
      },
      backdropBlur: {
        'green': 'blur(16px)'
      }
    }
  },
  plugins: [
    // Custom plugin for green-themed components
    function({ addComponents, theme }) {
      addComponents({
        '.btn-primary': {
          backgroundColor: theme('colors.primary.600'),
          color: theme('colors.white'),
          fontWeight: theme('fontWeight.semibold'),
          padding: `${theme('spacing.2')} ${theme('spacing.4')}`,
          borderRadius: theme('borderRadius.lg'),
          transition: 'all 0.2s ease-in-out',
          boxShadow: theme('boxShadow.green-sm'),
          '&:hover': {
            backgroundColor: theme('colors.primary.700'),
            boxShadow: theme('boxShadow.green-md'),
            transform: 'translateY(-1px)'
          },
          '&:active': {
            backgroundColor: theme('colors.primary.800'),
            transform: 'translateY(0)'
          }
        },
        '.btn-secondary': {
          backgroundColor: theme('colors.primary.100'),
          color: theme('colors.primary.800'),
          fontWeight: theme('fontWeight.semibold'),
          padding: `${theme('spacing.2')} ${theme('spacing.4')}`,
          borderRadius: theme('borderRadius.lg'),
          border: `1px solid ${theme('colors.primary.300')}`,
          transition: 'all 0.2s ease-in-out',
          '&:hover': {
            backgroundColor: theme('colors.primary.200'),
            borderColor: theme('colors.primary.400')
          }
        },
        '.card-green': {
          backgroundColor: theme('colors.white'),
          borderRadius: theme('borderRadius.xl'),
          padding: theme('spacing.6'),
          boxShadow: theme('boxShadow.green'),
          border: `1px solid ${theme('colors.primary.200')}`,
          transition: 'all 0.2s ease-in-out',
          '&:hover': {
            boxShadow: theme('boxShadow.green-md'),
            borderColor: theme('colors.primary.300')
          }
        },
        '.input-green': {
          width: '100%',
          padding: `${theme('spacing.3')} ${theme('spacing.4')}`,
          border: `1px solid ${theme('colors.primary.300')}`,
          borderRadius: theme('borderRadius.lg'),
          backgroundColor: theme('colors.white'),
          color: theme('colors.primary.900'),
          transition: 'all 0.2s ease-in-out',
          '&:focus': {
            outline: 'none',
            ringWidth: '2px',
            ringColor: theme('colors.primary.500'),
            borderColor: theme('colors.primary.500')
          },
          '&::placeholder': {
            color: theme('colors.primary.400')
          }
        },
        '.badge-green': {
          display: 'inline-flex',
          alignItems: 'center',
          padding: `${theme('spacing.1')} ${theme('spacing.3')}`,
          borderRadius: theme('borderRadius.full'),
          fontSize: theme('fontSize.xs'),
          fontWeight: theme('fontWeight.medium'),
          backgroundColor: theme('colors.primary.100'),
          color: theme('colors.primary.800')
        },
        '.nav-link-green': {
          color: theme('colors.primary.700'),
          padding: `${theme('spacing.2')} ${theme('spacing.3')}`,
          borderRadius: theme('borderRadius.md'),
          fontSize: theme('fontSize.sm'),
          fontWeight: theme('fontWeight.medium'),
          transition: 'all 0.15s ease-in-out',
          '&:hover': {
            backgroundColor: theme('colors.primary.100'),
            color: theme('colors.primary.900')
          },
          '&.active': {
            backgroundColor: theme('colors.primary.200'),
            color: theme('colors.primary.900')
          }
        },
        '.alert-success': {
          backgroundColor: theme('colors.primary.50'),
          border: `1px solid ${theme('colors.primary.200')}`,
          color: theme('colors.primary.800'),
          padding: theme('spacing.4'),
          borderRadius: theme('borderRadius.lg')
        },
        '.gradient-green': {
          backgroundImage: `linear-gradient(to right, ${theme('colors.primary.400')}, ${theme('colors.primary.600')})`
        },
        '.loading-spinner-green': {
          animation: 'spin 1s linear infinite',
          borderRadius: '50%',
          height: theme('spacing.8'),
          width: theme('spacing.8'),
          border: `2px solid ${theme('colors.primary.200')}`,
          borderTopColor: theme('colors.primary.600')
        }
      })
    }
  ]
}