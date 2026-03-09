// Configure React's act() environment for jsdom tests
globalThis.IS_REACT_ACT_ENVIRONMENT = true

// Initialize i18n with English for tests so assertions use English strings
import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import en from '@/locales/en.json'
import ja from '@/locales/ja.json'

void i18n.use(initReactI18next).init({
  resources: {
    en: { translation: en },
    ja: { translation: ja },
  },
  lng: 'en',
  fallbackLng: 'en',
  interpolation: { escapeValue: false },
  // Suppress i18next promotional output during tests
  partialBundledLanguages: true,
})
