import { ref, onMounted } from 'vue';
import { Environment } from '@/wailsjs/wailsjs/runtime/runtime';

const isMacOS = ref(false);
const isWindows = ref(false);
const isLinux = ref(false);
const platformDetected = ref(false);

export function usePlatform() {
  onMounted(async () => {
    if (platformDetected.value) {
      return; // Already detected
    }

    try {
      const env = await Environment();
      isMacOS.value = env.platform === 'darwin';
      isWindows.value = env.platform === 'windows';
      isLinux.value = env.platform === 'linux';
      platformDetected.value = true;
    } catch (error) {
      console.error('Failed to detect platform:', error);
      // Fallback to user agent detection
      const ua = navigator.userAgent.toLowerCase();
      isMacOS.value = ua.includes('mac');
      isWindows.value = ua.includes('win');
      isLinux.value = ua.includes('linux');
      platformDetected.value = true;
    }
  });

  return {
    isMacOS,
    isWindows,
    isLinux,
    platformDetected,
  };
}
