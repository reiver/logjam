import Logo from 'assets/images/logo.svg?react'

export function NotFound() {
  return (
    <div class="w-full h-screen dark:bg-black bg-white flex items-center justify-center">
      <div class="flex flex-col items-center justify-center gap-10 dark:text-gray-0 text-black">
        <Logo />
        <span class="dark:text-white text-black text-semi-bold-32">Oops! A 404 Error!</span>
        <span class="dark:text-white text-black text-medium-16">The page you are trying to reach is not found. Please edit the URL and try again.</span>
      </div>
    </div>
  )
}
export default NotFound
