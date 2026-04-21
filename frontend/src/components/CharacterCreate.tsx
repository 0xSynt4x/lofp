import { useState } from 'react'

const RACES = [
  {
    id: 1, name: '人类',
    stats: 'STR 30-100  AGI 30-100  QUI 30-100  CON 30-100  PER 30-100  WIL 40-110  EMP 30-100',
    desc: '破碎疆域的古老种族。人类自远古以来便已存在，凭借强大的意志力在各种技能上都能有所成就。赛博格技术最初是为人类开发的，因此他们可以使用任何此类装置。和其他种族相处的方式，就和他们彼此相处的方式一样。',
    ability: '可以使用所有赛博格植入物。',
  },
  {
    id: 2, name: '精灵',
    stats: 'STR 20-90  AGI 40-110  QUI 40-110  CON 1-70  PER 40-110  WIL 30-100  EMP 40-110',
    desc: '身材高挑纤细的类人生物，拥有美丽的容貌和尖耳朵。生活在林地中，敏捷迅速，感官敏锐。寿命可达数百年，对自然疾病具有极强的抵抗力。无忧无虑且超然物外，性情反复无常，是美与艺术的爱好者。',
    ability: 'CALL — 召唤林地生物为你服务（仅限野外）。',
  },
  {
    id: 3, name: '高地人',
    stats: 'STR 40-110  AGI 20-90  QUI 20-90  CON 50-120  PER 30-100  WIL 30-100  EMP 10-80',
    desc: '被称为石之子民的矮壮山地民族。强壮、坚韧，尤其对魔法具有抵抗力。在黑暗中视物清晰，熟练使用科技装置。他们的女性可以非常迷人，而且肯定没有胡子。',
    ability: 'BLEND — 融入洞穴/山地地形。5级时：MOLD — 将宝石塑造成更有价值的宝石。',
  },
  {
    id: 4, name: '狼人',
    stats: 'STR 30-100  AGI 40-110  QUI 40-110  CON 30-100  PER 40-110  WIL 30-100  EMP 30-100',
    desc: '中等身材的强壮类人生物，拥有明显的狼类特征。可以长途跋涉而不觉疲惫，夜间视力如同白昼。骄傲而高贵的种族，荣誉极为重要。称呼狼人为狼人是危险的举动。',
    ability: 'TRANSFORM — 变成巨狼形态——用爪牙战斗，长途跋涉。',
  },
  {
    id: 5, name: '穆格',
    stats: 'STR 40-110  AGI 30-100  QUI 30-100  CON 40-110  PER 40-110  WIL 20-90  EMP 20-90',
    desc: '天生具有竞争欲望的魁梧类人生物。强壮的四肢使他们成为熟练的攀爬者。在黑暗中视力几乎和白天一样好。忍不住要恶作剧，他们觉得自己非常有趣。分裂成不断互相争斗的氏族。',
    ability: 'FRENZY — 更凶猛地攻击；在生命值低于0时继续战斗直到死亡。',
  },
  {
    id: 6, name: '龙人',
    stats: 'STR 40-110  AGI 10-80  QUI 40-110  CON 40-110  PER 30-100  WIL 30-100  EMP 40-110',
    desc: '身覆坚硬鳞片的龙族后裔，拥有有力的下颚、爬行类的尾巴和巨大的蝙蝠状翅膀。必要时可以闪电般快速移动。他们的文化发展出独特的武斗流派和武器：武士刀、短刀、叉、锁镰、双节棍、棍棒、长柄刀、手里剑。绝不穿盔甲。',
    ability: 'FLY — 用翅膀飞行。独特的龙人武器风格，可与双武器战斗结合。',
  },
  {
    id: 7, name: '机械体',
    stats: 'STR 40-110  AGI 30-100  QUI 30-100  CON 40-110  PER 40-110  WIL 30-100  EMP 1-60',
    desc: '被注入了飘渺灵魂的机器。再生活组织覆盖着他们的机械身躯。是所有种族中最缺乏同理心的，使其成为最差的施法者，但他们可以随意开关情绪以最大化技能效果。务实且理性到令人厌烦的程度。',
    ability: 'EMOTE/UNEMOTE — 切换情绪状态以精确使用技能。',
  },
  {
    id: 8, name: '虚影',
    stats: 'STR n/a  AGI 30-100  QUI 50-120  CON 1-10  PER 30-100  WIL 30-100  EMP 30-100',
    desc: '不能完全存在于物质层面的飘渺幽灵。非魔法武器无法伤害他们，除非他们进行攻击——这会将他们拉入物质层面。他们无法操控重物。天生的心灵感应者，通常致力于学术追求。',
    ability: '对非魔法武器免疫（攻击时除外）。天生心灵感应。黑暗中视物。',
  },
]

interface Props {
  onCreated: (char: { firstName: string; lastName: string; race: number; gender: number }) => void
  onOpenManual?: () => void
}

export default function CharacterCreate({ onCreated, onOpenManual }: Props) {
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [race, setRace] = useState(1)
  const [gender, setGender] = useState(0)
  const [selectedRace, setSelectedRace] = useState(RACES[0])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!firstName.trim() || !lastName.trim()) return
    onCreated({ firstName: firstName.trim(), lastName: lastName.trim(), race, gender })
  }

  return (
    <div className="flex items-start sm:items-center justify-center h-full pt-4 px-4 pb-4 sm:p-8 overflow-y-auto">
      <div className="max-w-3xl w-full bg-[#111] border border-[#333] rounded-lg p-4 sm:p-8">
        <h2 className="text-amber-400 text-2xl font-mono mb-1 text-center">
          创建你的角色
        </h2>
        <p className="text-gray-500 text-sm font-mono mb-6 text-center">
          谨慎选择——你的种族和能力将决定你在破碎疆域中的命运
        </p>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Name */}
          <div className="bg-[#0a0a0a] border border-amber-900/50 rounded-lg p-3 mb-2">
            <p className="text-gray-400 text-xs font-mono leading-relaxed">
              Legends 是一个角色扮演游戏——请选择适合奇幻设定的名字。
              避免现代名字、流行文化引用或玩笑名字.{' '}
              <button type="button" onClick={onOpenManual} className="text-amber-500 hover:text-amber-400 underline cursor-pointer">
                了解更多关于角色扮演的信息 &rarr;
              </button>
            </p>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block text-gray-400 text-sm font-mono mb-1">名</label>
              <input
                type="text"
                value={firstName}
                onChange={e => setFirstName(e.target.value)}
                maxLength={15}
                className="w-full bg-[#0a0a0a] border border-[#444] rounded px-3 py-2 text-gray-200 font-mono focus:border-amber-500 focus:outline-none"
                placeholder="巴尔萨泽"
                autoFocus
              />
            </div>
            <div>
              <label className="block text-gray-400 text-sm font-mono mb-1">姓</label>
              <input
                type="text"
                value={lastName}
                onChange={e => setLastName(e.target.value)}
                maxLength={15}
                className="w-full bg-[#0a0a0a] border border-[#444] rounded px-3 py-2 text-gray-200 font-mono focus:border-amber-500 focus:outline-none"
                placeholder="辛瓦尔"
              />
              <p className="text-gray-600 text-[10px] font-mono mt-1">姓氏允许使用连字符和重音符号</p>
            </div>
          </div>

          {/* Gender */}
          <div>
            <label className="block text-gray-400 text-sm font-mono mb-1">性别</label>
            <div className="flex gap-4">
              {[{ v: 0, l: '男性' }, { v: 1, l: '女性' }].map(g => (
                <button
                  key={g.v}
                  type="button"
                  onClick={() => setGender(g.v)}
                  className={`px-6 py-2.5 min-h-[44px] rounded font-mono text-sm transition-colors ${gender === g.v ? 'bg-amber-700 text-white' : 'bg-[#1a1a1a] text-gray-400 border border-[#444] hover:border-amber-600'}`}
                >
                  {g.l}
                </button>
              ))}
            </div>
          </div>

          {/* Race selection */}
          <div>
            <label className="block text-gray-400 text-sm font-mono mb-2">种族</label>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 mb-3">
              {RACES.map(r => (
                <button
                  key={r.id}
                  type="button"
                  onClick={() => { setRace(r.id); setSelectedRace(r) }}
                  className={`px-2 py-3 min-h-[44px] rounded font-mono text-sm transition-colors ${race === r.id ? 'bg-amber-700 text-white border border-amber-600' : 'bg-[#1a1a1a] text-gray-400 border border-[#444] hover:border-amber-600'}`}
                >
                  {r.name}
                </button>
              ))}
            </div>

            {/* Race detail */}
            <div className="bg-[#0a0a0a] border border-[#333] rounded-lg p-4 space-y-3">
              <div className="flex items-center gap-3">
                <h3 className="text-amber-400 font-mono text-lg font-bold">{selectedRace.name}</h3>
              </div>
              <p className="text-gray-300 font-mono text-xs leading-relaxed">{selectedRace.desc}</p>
              <div className="bg-[#111] border border-[#2a2a2a] rounded p-2">
                <p className="text-green-400 font-mono text-xs">
                  <span className="text-gray-500">Ability:</span> {selectedRace.ability}
                </p>
              </div>
              <div className="bg-[#111] border border-[#2a2a2a] rounded p-2">
                <p className="text-cyan-400 font-mono text-[10px] tracking-wider">{selectedRace.stats}</p>
              </div>
            </div>
          </div>

          <button
            type="submit"
            disabled={!firstName.trim() || !lastName.trim()}
            className="w-full py-3 bg-amber-700 hover:bg-amber-600 disabled:bg-gray-700 disabled:text-gray-500 text-white font-mono rounded text-lg transition-colors cursor-pointer"
          >
            进入破碎疆域
          </button>
        </form>
      </div>
    </div>
  )
}
