import { useState } from 'react'
import './App.css'
import ChartComponent from './components/ChartComponent'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-blue-600 text-white p-6">
        <div className="max-w-6xl mx-auto">
          <h1 className="text-3xl font-bold">StrikePad Frontend</h1>
          <p className="text-blue-100 mt-2">React + TypeScript + Tailwind + D3.js</p>
        </div>
      </header>
      
      <main className="max-w-6xl mx-auto p-6">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-white p-6 rounded-lg shadow-md">
            <h2 className="text-xl font-semibold mb-4">Counter Demo</h2>
            <div className="space-y-4">
              <button 
                onClick={() => setCount((count) => count + 1)}
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded transition-colors"
              >
                Count is {count}
              </button>
              <p className="text-gray-600">
                Edit <code className="bg-gray-100 px-2 py-1 rounded">src/App.tsx</code> and save to test HMR
              </p>
            </div>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow-md">
            <h2 className="text-xl font-semibold mb-4">D3.js Chart Demo</h2>
            <ChartComponent />
          </div>
        </div>
      </main>
    </div>
  )
}

export default App