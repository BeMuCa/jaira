import type { Lang } from './lang.ts'

// The cockpit's words for a person, one table per language lang.ts knows.

type Texts = {
  waiting: string
  accept: string
  giveBack: string
  explain: string
  discuss: string
  more: (n: number) => string
  reasonLabel: string
  reasonPlaceholder: string
  accepted: (id: string) => string
  returnedTo: (id: string, who: 'old' | 'new' | 'none') => string
  failed: (what: string, why: string) => string
  needsHerdr: string
  opened: (id: string) => string
  context: (k: number) => string
  cacheLeft: (min: number) => string
  cacheCold: string
  hintCompact: string
  hintNew: string
  guardQuestion: (k: number, min: number, usd: string) => string
  guardSend: string
  guardCompact: string
  guardClear: string
  guardCancel: string
  guardHeader: string
  threadsOpen: (questions: number, stale: number) => string
  threadsNone: string
  threadsTitle: string
  status: (worker: number, review: number, person: number) => string
  usage: string
  paneTitle: (id: string) => string
  paneReports: string
  paneUsage: string
  paneNothing: string
  usageText: string
}

const T: Record<Lang, Texts> = {
  en: {
    waiting: 'Waiting for you',
    accept: 'Accept',
    giveBack: 'Give back',
    explain: 'Explain',
    discuss: 'Discuss',
    more: n => `+${n} more`,
    reasonLabel: 'What is missing?',
    reasonPlaceholder: 'goes onto the ticket as a note',
    accepted: id => `${id} accepted: done`,
    returnedTo: (id, who) =>
      who === 'none' ? `${id} back in in-progress; no worker started (no Herdr)` : `${id} given back to ${who === 'old' ? 'the worker that carried it' : 'a new worker'}`,
    failed: (what, why) => `${what} failed: ${why}`,
    needsHerdr: 'A chat of its own needs Herdr (HERDR_ENV=1).',
    opened: id => `${id}: opened in a Herdr tab of its own`,
    context: k => `context ${k}k`,
    cacheLeft: min => `cache ${min} min left`,
    cacheCold: 'cache cold',
    hintCompact: '/compact advised',
    hintNew: 'start a new session',
    guardQuestion: (k, min, usd) => `The prompt cache went cold ${min} min ago. Sending writes ${k}k tokens to it again (about ${usd}). What now?`,
    guardSend: 'Send',
    guardCompact: 'Compact first',
    guardClear: 'Start over (/clear)',
    guardCancel: 'Cancel',
    guardHeader: 'Cold cache',
    threadsOpen: (q, s) => `threads: ${q} open question(s), ${s} left lying`,
    threadsNone: 'No conversation threads yet.',
    threadsTitle: 'Conversation threads (? open question to you, ~ left lying, -> ticket):',
    status: (w, r, p) => `${w} worker · ${r} review · ${p} waiting for you`,
    usage: 'session',
    paneTitle: id => `Ticket ${id}`,
    paneReports: 'Reports',
    paneUsage: 'Usage per agent',
    paneNothing: 'nothing recorded yet',
    usageText:
      'Usage: /cockpit <subcommand>\n\n  report [id]   full text of the agent reports shortened for a ticket, newest first; without an id, the tickets that have some\n  threads       the conversation threads of this session and where each stands\n  ticket <id>   a pane with the ticket\'s reports and what its agents used',
  },
  de: {
    waiting: 'Wartet auf dich',
    accept: 'Annehmen',
    giveBack: 'Zurückgeben',
    explain: 'Erklären',
    discuss: 'Besprechen',
    more: n => `+${n} weitere`,
    reasonLabel: 'Was fehlt?',
    reasonPlaceholder: 'kommt als Notiz ans Ticket',
    accepted: id => `${id} angenommen: done`,
    returnedTo: (id, who) =>
      who === 'none' ? `${id} wieder in in-progress; kein Worker gestartet (kein Herdr)` : `${id} zurückgegeben an ${who === 'old' ? 'den Worker, der es getragen hat' : 'einen neuen Worker'}`,
    failed: (what, why) => `${what} fehlgeschlagen: ${why}`,
    needsHerdr: 'Ein eigener Chat braucht Herdr (HERDR_ENV=1).',
    opened: id => `${id}: in einem eigenen Herdr-Tab geöffnet`,
    context: k => `Kontext ${k}k`,
    cacheLeft: min => `Cache noch ${min} min`,
    cacheCold: 'Cache kalt',
    hintCompact: '/compact empfohlen',
    hintNew: 'neue Sitzung beginnen',
    guardQuestion: (k, min, usd) => `Der Prompt-Cache ist seit ${min} min kalt. Senden schreibt ${k}k Tokens neu hinein (etwa ${usd}). Was nun?`,
    guardSend: 'Senden',
    guardCompact: 'Erst komprimieren',
    guardClear: 'Neu anfangen (/clear)',
    guardCancel: 'Abbrechen',
    guardHeader: 'Cache kalt',
    threadsOpen: (q, s) => `Fäden: ${q} offene Frage(n), ${s} liegen geblieben`,
    threadsNone: 'Noch keine Gesprächsfäden.',
    threadsTitle: 'Gesprächsfäden (? offene Frage an dich, ~ liegen geblieben, -> Ticket):',
    status: (w, r, p) => `${w} worker · ${r} review · ${p} wartet auf dich`,
    usage: 'Sitzung',
    paneTitle: id => `Ticket ${id}`,
    paneReports: 'Berichte',
    paneUsage: 'Verbrauch je Agent',
    paneNothing: 'noch nichts erfasst',
    usageText:
      'Aufruf: /cockpit <unterbefehl>\n\n  report [id]   voller Text der gekürzten Agentenberichte zu einem Ticket, neueste zuerst; ohne id die Tickets, zu denen es welche gibt\n  threads       die Gesprächsfäden dieser Sitzung und wo jeder steht\n  ticket <id>   ein Pane mit den Berichten des Tickets und dem Verbrauch seiner Agenten',
  },
  ru: {
    waiting: 'Ждёт тебя',
    accept: 'Принять',
    giveBack: 'Вернуть',
    explain: 'Объяснить',
    discuss: 'Обсудить',
    more: n => `ещё ${n}`,
    reasonLabel: 'Чего не хватает?',
    reasonPlaceholder: 'попадёт в тикет заметкой',
    accepted: id => `${id} принят: done`,
    returnedTo: (id, who) =>
      who === 'none' ? `${id} снова в in-progress; воркер не запущен (нет Herdr)` : `${id} возвращён ${who === 'old' ? 'воркеру, который его вёл' : 'новому воркеру'}`,
    failed: (what, why) => `${what} не удалось: ${why}`,
    needsHerdr: 'Для отдельного чата нужен Herdr (HERDR_ENV=1).',
    opened: id => `${id}: открыт в отдельной вкладке Herdr`,
    context: k => `контекст ${k}k`,
    cacheLeft: min => `кэш ещё ${min} мин`,
    cacheCold: 'кэш остыл',
    hintCompact: 'советую /compact',
    hintNew: 'начни новую сессию',
    guardQuestion: (k, min, usd) => `Кэш промпта остыл ${min} мин назад. Отправка заново запишет в него ${k}k токенов (примерно ${usd}). Что делаем?`,
    guardSend: 'Отправить',
    guardCompact: 'Сначала сжать',
    guardClear: 'Начать заново (/clear)',
    guardCancel: 'Отмена',
    guardHeader: 'Кэш остыл',
    threadsOpen: (q, s) => `нити: открытых вопросов ${q}, заброшено ${s}`,
    threadsNone: 'Нитей разговора пока нет.',
    threadsTitle: 'Нити разговора (? открытый вопрос к тебе, ~ заброшено, -> тикет):',
    status: (w, r, p) => `${w} worker · ${r} review · ${p} ждёт тебя`,
    usage: 'сессия',
    paneTitle: id => `Тикет ${id}`,
    paneReports: 'Отчёты',
    paneUsage: 'Расход по агентам',
    paneNothing: 'пока ничего не записано',
    usageText:
      'Вызов: /cockpit <подкоманда>\n\n  report [id]   полный текст сокращённых отчётов агентов по тикету, свежие сверху; без id — тикеты, по которым они есть\n  threads       нити разговора этой сессии и где каждая стоит\n  ticket <id>   панель с отчётами тикета и расходом его агентов',
  },
}

export function texts(lang: Lang): Texts {
  return T[lang]
}
