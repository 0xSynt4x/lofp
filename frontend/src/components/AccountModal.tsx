import { useState } from 'react'
import { useAuth } from '../App'

interface Props {
  onClose: () => void
}

export default function AccountModal({ onClose }: Props) {
  const { user, logout } = useAuth()
  const [tab, setTab] = useState<'info' | 'name' | 'password' | 'verify'>('info')
  const [newName, setNewName] = useState(user?.account?.name || '')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [verifyCode, setVerifyCode] = useState('')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const headers = (): Record<string, string> => {
    const h: Record<string, string> = { 'Content-Type': 'application/json' }
    if (user?.token) h['Authorization'] = `Bearer ${user.token}`
    return h
  }

  const handleUpdateName = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true); setError(''); setMessage('')
    try {
      const r = await fetch('/api/auth/me/name', {
        method: 'PUT', headers: headers(),
        body: JSON.stringify({ name: newName }),
      })
      if (r.ok) {
        setMessage('Display name updated!')
        // Update local storage
        const stored = localStorage.getItem('lofp_auth')
        if (stored) {
          const parsed = JSON.parse(stored)
          parsed.account.name = newName
          localStorage.setItem('lofp_auth', JSON.stringify(parsed))
        }
      } else {
        const d = await r.json().catch(() => null)
        setError(d?.error || 'Failed to update name')
      }
    } catch { setError('Network error') }
    setSubmitting(false)
  }

  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    if (password !== confirmPassword) { setError('Passwords do not match'); return }
    setSubmitting(true); setError(''); setMessage('')
    try {
      const r = await fetch('/api/auth/me/password', {
        method: 'PUT', headers: headers(),
        body: JSON.stringify({ password }),
      })
      if (r.ok) {
        setMessage('Password updated!')
        setPassword(''); setConfirmPassword('')
      } else {
        const d = await r.json().catch(() => null)
        setError(d?.error || 'Failed to update password')
      }
    } catch { setError('Network error') }
    setSubmitting(false)
  }

  const handleResendVerification = async () => {
    setSubmitting(true); setError(''); setMessage('')
    try {
      const r = await fetch('/api/auth/resend-verification', {
        method: 'POST', headers: headers(),
        body: JSON.stringify({ email: user?.account?.email }),
      })
      if (r.ok) {
        setMessage('Verification email sent! Check your inbox.')
      } else {
        const d = await r.json().catch(() => null)
        setError(d?.error || 'Failed to send verification email')
      }
    } catch { setError('Network error') }
    setSubmitting(false)
  }

  const handleVerifyCode = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true); setError(''); setMessage('')
    try {
      const r = await fetch('/api/auth/verify-code', {
        method: 'POST', headers: headers(),
        body: JSON.stringify({ code: verifyCode }),
      })
      if (r.ok) {
        setMessage('Email verified!')
        // Update local storage
        const stored = localStorage.getItem('lofp_auth')
        if (stored) {
          const parsed = JSON.parse(stored)
          parsed.account.emailVerified = true
          localStorage.setItem('lofp_auth', JSON.stringify(parsed))
        }
        // Reload page to refresh auth state
        setTimeout(() => window.location.reload(), 1000)
      } else {
        const d = await r.json().catch(() => null)
        setError(d?.error || 'Invalid verification code')
      }
    } catch { setError('Network error') }
    setSubmitting(false)
  }

  const isVerified = user?.account?.emailVerified !== false

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-[#1a1a1a] border border-[#444] rounded-lg p-6 max-w-md w-full mx-4" onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-amber-400 font-mono font-bold text-lg">账户设置</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-300 text-lg">&times;</button>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 mb-4 border-b border-[#333] pb-2">
          <button onClick={() => { setTab('info'); setError(''); setMessage('') }}
            className={`px-3 py-1 text-xs font-mono rounded-t ${tab === 'info' ? 'bg-[#333] text-amber-400' : 'text-gray-500 hover:text-gray-300'}`}>
            信息
          </button>
          <button onClick={() => { setTab('name'); setError(''); setMessage('') }}
            className={`px-3 py-1 text-xs font-mono rounded-t ${tab === 'name' ? 'bg-[#333] text-amber-400' : 'text-gray-500 hover:text-gray-300'}`}>
            名称
          </button>
          <button onClick={() => { setTab('password'); setError(''); setMessage('') }}
            className={`px-3 py-1 text-xs font-mono rounded-t ${tab === 'password' ? 'bg-[#333] text-amber-400' : 'text-gray-500 hover:text-gray-300'}`}>
            密码
          </button>
          {!isVerified && (
            <button onClick={() => { setTab('verify'); setError(''); setMessage('') }}
              className={`px-3 py-1 text-xs font-mono rounded-t ${tab === 'verify' ? 'bg-[#333] text-amber-400' : 'text-yellow-500 hover:text-yellow-300 animate-pulse'}`}>
              验证邮箱
            </button>
          )}
        </div>

        {/* Tab content */}
        {tab === 'info' && (
          <div className="space-y-2 font-mono text-sm">
            <div><span className="text-gray-500">邮箱:</span> <span className="text-gray-300">{user?.account?.email}</span></div>
            <div><span className="text-gray-500">名称:</span> <span className="text-gray-300">{user?.account?.name}</span></div>
            <div>
              <span className="text-gray-500">邮箱已验证:</span>{' '}
              {isVerified
                ? <span className="text-green-400">是</span>
                : <span className="text-yellow-400">否 — <button onClick={() => setTab('verify')} className="underline hover:text-yellow-300">立即验证</button></span>
              }
            </div>
            <div className="flex items-center gap-2 mt-2">
              <img src={user?.account?.picture || '/default-avatar.svg'} alt="" className="w-8 h-8 rounded-full" />
              {user?.account?.picture
                ? <span className="text-gray-500 text-xs">已关联 Google</span>
                : <span className="text-gray-500 text-xs">邮箱/密码账户</span>
              }
            </div>
            <div className="pt-3 mt-3 border-t border-[#333]">
              <button onClick={logout} className="text-red-400 hover:text-red-300 text-xs font-mono">
                登出
              </button>
            </div>
          </div>
        )}

        {tab === 'name' && (
          <form onSubmit={handleUpdateName} className="space-y-3">
            <input type="text" value={newName} onChange={e => setNewName(e.target.value)}
              placeholder="显示名称"
              className="w-full px-3 py-2 bg-[#111] border border-[#444] rounded font-mono text-sm text-gray-200 focus:border-amber-600 focus:outline-none" />
            {error && <p className="text-red-400 font-mono text-xs">{error}</p>}
            {message && <p className="text-green-400 font-mono text-xs">{message}</p>}
            <button type="submit" disabled={submitting}
              className="w-full py-2 bg-amber-700 hover:bg-amber-600 text-white font-mono text-sm rounded disabled:opacity-50">
              {submitting ? '更新中...' : '更新名称'}
            </button>
          </form>
        )}

        {tab === 'password' && (
          <form onSubmit={handleUpdatePassword} className="space-y-3">
            {user?.account?.picture && (
              <p className="text-gray-500 font-mono text-xs">
                设置密码以启用 telnet 或 SSH MUD 客户端登录。
              </p>
            )}
            <input type="password" value={password} onChange={e => setPassword(e.target.value)}
              placeholder="新密码 (10+ 字符, 大小写混合, 数字, 特殊字符)"
              className="w-full px-3 py-2 bg-[#111] border border-[#444] rounded font-mono text-sm text-gray-200 focus:border-amber-600 focus:outline-none" />
            <input type="password" value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)}
              placeholder="确认新密码"
              className="w-full px-3 py-2 bg-[#111] border border-[#444] rounded font-mono text-sm text-gray-200 focus:border-amber-600 focus:outline-none" />
            {error && <p className="text-red-400 font-mono text-xs">{error}</p>}
            {message && <p className="text-green-400 font-mono text-xs">{message}</p>}
            <button type="submit" disabled={submitting}
              className="w-full py-2 bg-amber-700 hover:bg-amber-600 text-white font-mono text-sm rounded disabled:opacity-50">
              {submitting ? '更新中...' : '更新密码'}
            </button>
          </form>
        )}

        {tab === 'verify' && (
          <div className="space-y-3">
            <p className="text-gray-400 font-mono text-sm">
              输入邮件中的验证码，或点击下方重新发送。
            </p>
            <form onSubmit={handleVerifyCode} className="space-y-3">
              <input type="text" value={verifyCode} onChange={e => setVerifyCode(e.target.value.toUpperCase())}
                placeholder="验证码 (例如 ABCD1234)"
                className="w-full px-3 py-2 bg-[#111] border border-[#444] rounded font-mono text-sm text-gray-200 focus:border-amber-600 focus:outline-none tracking-widest text-center text-lg"
                maxLength={8} autoFocus />
              {error && <p className="text-red-400 font-mono text-xs">{error}</p>}
              {message && <p className="text-green-400 font-mono text-xs">{message}</p>}
              <button type="submit" disabled={submitting || verifyCode.length < 8}
                className="w-full py-2 bg-amber-700 hover:bg-amber-600 text-white font-mono text-sm rounded disabled:opacity-50">
                {submitting ? '验证中...' : '验证'}
              </button>
            </form>
            <button onClick={handleResendVerification} disabled={submitting}
              className="w-full py-2 bg-[#222] hover:bg-[#333] text-gray-400 font-mono text-xs rounded border border-[#444] disabled:opacity-50">
              重新发送验证邮件
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
