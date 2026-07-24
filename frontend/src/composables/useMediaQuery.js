import { ref, onMounted, onBeforeUnmount } from 'vue'

export function useMediaQuery(query) {
  const matches = ref(false)
  let mql = null
  const onChange = (e) => { matches.value = e.matches }

  onMounted(() => {
    mql = window.matchMedia(query)
    matches.value = mql.matches
    mql.addEventListener('change', onChange)
  })
  onBeforeUnmount(() => mql?.removeEventListener('change', onChange))

  return matches
}

/** 移动端断点（<768px） */
export const useIsMobile = () => useMediaQuery('(max-width: 767px)')
