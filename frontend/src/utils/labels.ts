export const kindNames: Record<string, string> = {
  food: 'Продукты',
  clothes: 'Вещи',
  furniture: 'Мебель',
  stationery: 'Канцелярия',
  gifts: 'Подарки',
  consultation: 'Консультация',
  other: 'Другое'
}
export const profileStatusNames: Record<string, string> = {
  draft: 'Черновик',
  sent: 'Отправлена',
  checked: 'Проверена'
}
export const motherStatusNames: Record<string, string> = {
  review: 'На рассмотрении',
  active: 'На сопровождении',
  inactive: 'Неактивна'
}
export const requestStatusNames: Record<string, string> = {
  new: 'Новая',
  in_work: 'В работе',
  issued: 'Выдана',
  rejected: 'Отклонена'
}
export const participantStatusNames: Record<string, string> = {
  registered: 'Записана',
  came: 'Пришла',
  absent: 'Не пришла'
}
export const docTypeNames: Record<string, string> = {
  questionnaire: 'Анкета',
  consent: 'Согласие на обработку данных',
  responsibility: 'Заявление об ответственности',
  other: 'Другое'
}
export const donorKindNames: Record<string, string> = {
  person: 'Человек',
  company: 'Организация'
}
export function label(names: Record<string, string>, key: string) {
  return names[key] || key
}
export function sumText(kind: string, sum: number) {
  if (kind === 'food') return sum + ' ₽'
  if (kind === 'consultation') return String(sum) + ' визитов'
  return sum + ' шт.'
}
