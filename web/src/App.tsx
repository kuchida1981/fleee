import { Button } from '@/components/ui/button'

function App() {
  return (
    <div className="min-h-screen bg-neutral-50 dark:bg-neutral-900 text-neutral-900 dark:text-neutral-50 font-sans">
      <header className="border-b border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-950 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-2xl font-bold tracking-tight text-neutral-950 dark:text-neutral-50">fleee</span>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <h1 className="text-3xl font-bold tracking-tight text-neutral-950 dark:text-neutral-50">
              勘定科目管理
            </h1>
          </div>
          <div className="border border-dashed border-neutral-200 dark:border-neutral-800 rounded-lg p-12 text-center text-neutral-500 dark:text-neutral-400">
            <p className="mb-4">UIコンポーネントは次のステップで実装します。</p>
            <Button onClick={() => alert('Button Clicked!')}>動作確認ボタン</Button>
          </div>
        </div>
      </main>
    </div>
  )
}

export default App

