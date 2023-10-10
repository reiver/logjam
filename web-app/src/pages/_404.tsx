import Logo from 'assets/images/Logo + Type.svg?react'

export function NotFound() {
  return (
    <div class="w-full h-screen bg-primary flex items-center justify-center">
      <div class="flex flex-col items-center justify-center gap-10">
        <Logo />
        <span class="text-medium-12">Page Not Found!</span>
      </div>
    </div>
  )
}
export default NotFound
