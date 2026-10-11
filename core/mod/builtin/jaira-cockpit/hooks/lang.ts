// Every text the cockpit shows a person, in the language the person's Claude
// Code `language` setting names. Ticket ids, lane names and commands stay as
// they are in every language: they are what the person types back.

export type Lang = 'en' | 'de' | 'ru'

// `language` is free text ("Russian", "русский", "de"); read it by its start.
const SPELLINGS: ReadonlyArray<[Lang, RegExp]> = [
  ['de', /^\s*(de\b|deu|ger|deutsch)/i],
  ['ru', /^\s*(ru\b|rus|русск)/i],
  ['en', /^\s*(en\b|eng)/i],
]

export function langOf(setting: unknown): Lang {
  if (typeof setting !== 'string') return 'en'
  return SPELLINGS.find(([, re]) => re.test(setting))?.[0] ?? 'en'
}


// The noise filter's own words: the digest line and /cockpit report. Kept
// here rather than with the cockpit's other texts so the filter stands alone.
type FilterWords = {
  // Leads every digest, so the model never takes the line for the agent's
  // own words.
  marker: string
  question: string
  fullReport: string
  none: string
  noneFor: (id: string) => string
  list: string
  hint: string
  count: (n: number) => string
  last: string
  of: string
}

const FILTER: Record<Lang, FilterWords> = {
  en: {
    marker: '[shortened by jaira-cockpit]',
    question: 'question',
    fullReport: 'full report',
    none: 'No agent reports kept for this project yet.',
    noneFor: id => `No reports kept for ${id} in this project.`,
    list: 'Tickets with kept reports in this project, most recent first:',
    hint: '/cockpit report <id> shows one.',
    count: n => `${n} report(s)`,
    last: 'last',
    of: 'of',
  },
  de: {
    marker: '[gekürzt von jaira-cockpit]',
    question: 'Frage',
    fullReport: 'voller Bericht',
    none: 'Zu diesem Projekt sind noch keine Agentenberichte aufbewahrt.',
    noneFor: id => `Zu ${id} sind in diesem Projekt keine Berichte aufbewahrt.`,
    list: 'Tickets mit aufbewahrten Berichten in diesem Projekt, neueste zuerst:',
    hint: '/cockpit report <id> zeigt einen.',
    count: n => `${n} Bericht(e)`,
    last: 'zuletzt',
    of: 'von',
  },
  ru: {
    marker: '[сокращено jaira-cockpit]',
    question: 'вопрос',
    fullReport: 'полный отчёт',
    none: 'По этому проекту отчётов агентов пока нет.',
    noneFor: id => `По ${id} в этом проекте отчётов нет.`,
    list: 'Тикеты с сохранёнными отчётами в этом проекте, свежие сверху:',
    hint: '/cockpit report <id> покажет один.',
    count: n => `отчётов: ${n}`,
    last: 'последний',
    of: 'из',
  },
}

export function filterWords(lang: Lang): FilterWords {
  return FILTER[lang]
}
