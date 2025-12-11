/**
 * Media proxy utilities for handling anti-hotlinking and caching
 */

// Cache for media cache enabled setting to avoid repeated API calls
let mediaCacheEnabledCache: boolean | null = null;
let mediaCachePromise: Promise<boolean> | null = null;

/**
 * Convert a media URL to use the proxy endpoint
 * @param url Original media URL
 * @param referer Optional referer URL for anti-hotlinking
 * @returns Proxied URL
 */
export function getProxiedMediaUrl(url: string, referer?: string): string {
  if (!url) return '';

  // Don't proxy data URLs or blob URLs
  if (url.startsWith('data:') || url.startsWith('blob:')) {
    return url;
  }

  // Don't proxy local URLs
  if (
    url.startsWith('/') ||
    url.startsWith('http://localhost') ||
    url.startsWith('http://127.0.0.1')
  ) {
    return url;
  }

  // Build proxy URL
  const params = new URLSearchParams();
  params.set('url', url);
  if (referer) {
    params.set('referer', referer);
  }

  return `/api/media/proxy?${params.toString()}`;
}

/**
 * Check if media caching is enabled (with caching to avoid repeated API calls)
 * @returns Promise<boolean>
 */
export async function isMediaCacheEnabled(): Promise<boolean> {
  // Return cached value if available
  if (mediaCacheEnabledCache !== null) {
    return mediaCacheEnabledCache;
  }

  // If a request is already in flight, wait for it
  if (mediaCachePromise) {
    return mediaCachePromise;
  }

  // Start a new request
  mediaCachePromise = (async () => {
    try {
      const response = await fetch('/api/settings');
      if (response.ok) {
        const settings = await response.json();
        mediaCacheEnabledCache =
          settings.media_cache_enabled === 'true' || settings.media_cache_enabled === true;
        return mediaCacheEnabledCache;
      }
    } catch (error) {
      console.error('Failed to check media cache status:', error);
    }
    mediaCacheEnabledCache = false;
    return false;
  })();

  const result = await mediaCachePromise;
  mediaCachePromise = null; // Clear the promise after completion
  return result;
}

/**
 * Clear the media cache enabled cache (call this when settings change)
 */
export function clearMediaCacheEnabledCache(): void {
  mediaCacheEnabledCache = null;
}

/**
 * Process HTML content to proxy image URLs
 * @param html HTML content
 * @param referer Optional referer URL
 * @returns HTML with proxied image URLs
 */
export function proxyImagesInHtml(html: string, referer?: string): string {
  if (!html) return html;

  // Replace img src attributes - handles double quotes, single quotes, and unquoted values
  return html.replace(/<img([^>]+)src\s*=\s*(['"]?)([^"'\s>]+)\2/gi, (match, attrs, quote, src) => {
    const proxiedUrl = getProxiedMediaUrl(src, referer);
    // Always output with double quotes for consistency
    return `<img${attrs}src="${proxiedUrl}"`;
  });
}
