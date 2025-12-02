/**
 * Centralized default values for settings
 * Modify these values to change defaults across the application
 */

export const settingsDefaults = {
  // General settings
  update_interval: 10,
  language: 'en-US',
  theme: 'auto',
  default_view_mode: 'original',
  startup_on_boot: false,
  show_hidden_articles: false,

  // Translation settings
  translation_enabled: false,
  target_language: 'zh',
  translation_provider: 'google',
  deepl_api_key: '',

  // Summary settings
  summary_enabled: false,
  summary_length: 'medium',

  // Cleanup settings
  auto_cleanup_enabled: false,
  max_cache_size_mb: 20,
  max_article_age_days: 30,

  // Other settings
  shortcuts: '',
  rules: '',
  last_article_update: '',
} as const;

// Type for the defaults object
export type SettingsDefaults = typeof settingsDefaults;
