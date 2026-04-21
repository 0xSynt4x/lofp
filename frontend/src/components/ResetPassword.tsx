import { useState } from 'react'

interface Props {
  onBack: () => void
}

export default function ResetPassword({ onBack }: Props) {
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [status, setStatus] = useState<'form' | 'success' | 'error'>('form')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const params = new URLSearchParams(window.location.search)
  const token = params.get('token')

  if (!token) {
    return (
      <div className="flex items-center justify-center h-full p-8">
        <div className="max-w-md w-full text-center">
          <p className="text-red-400 font-mono text-lg mb-4">缺少重置令牌。</p>
          <button onClick={onBack} className="px-6 py-2 bg-[#333] hover:bg-[#444] text-gray-300 font-mono rounded">
            返回菜单
          </button>
        </div>
      </div>
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (password !== confirm) {
      setError('两次输入的密码不一致。')
      return
    }
    setSubmitting(true)
    setError('')
    try {
      const resp = await fetch('/api/auth/reset-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, password }),
      })
      if (resp.ok) {
        setStatus('success')
      } else {
        const data = await resp.json().catch(() => null)
        setError(data?.error || '重置失败。')
        setStatus('error')
      }
    } catch {
      setError('网络错误。')
      setStatus('error')
    }
    setSubmitting(false)
  }

  return (
    <div className="flex items-center justify-center h-full p-8">
      <div className="max-w-sm w-full">
        {status === 'form' && (
          <div>
            <h2 className="text-amber-400 font-mono font-bold text-lg mb-4 text-center">设置新密码</h2>
            <form onSubmit={handleSubmit} className="space-y-3">
              <input
                type="password" placeholder="新密码 (10+字符, 大小写混合, 数字, 特殊字符)" value={password}
                onChange={e => setPassword(e.target.value)}
                className="w-full px-3 py-2 bg-[#111] border border-[#444] rounded font-mono text-sm text-gray-200 focus:border-amber-600 focus:outline-none"
                autoFocus
              />
              <input
                type="password" placeholder="确认密码" value={confirm}
                onChange={e => setConfirm(e.target.value)}
                className="w-full px-3 py-2 bg-[#111] border border-[#444] rounded font-mono text-sm text-gray-200 focus:border-amber-600 focus:outline-none"
              />
              {error && <p className="text-red-400 font-mono text-xs">{error}</p>}
              <button type="submit" disabled={submitting}
                className="w-full py-2 bg-amber-700 hover:bg-amber-600 text-white font-mono text-sm rounded disabled:opacity-50 transition-colors">
                {submitting ? '重置中...' : '重置密码'}
              </button>
            </form>
          </div>
        )}
        {status === 'success' && (
          <div className="text-center">
            <p className="text-green-400 font-mono text-lg mb-4">密码已更新！</p>
            <button onClick={onBack} className="px-6 py-2 bg-amber-700 hover:bg-amber-600 text-white font-mono rounded">
              登录
            </button>
          </div>
        )}
        {status === 'error' && (
          <div className="text-center">
            <p className="text-red-400 font-mono text-lg mb-4">重置失败</p>
            <p className="text-gray-400 font-mono text-sm mb-6">{error}</p>
            <button onClick={onBack} className="px-6 py-2 bg-[#333] hover:bg-[#444] text-gray-300 font-mono rounded">
              返回菜单
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
