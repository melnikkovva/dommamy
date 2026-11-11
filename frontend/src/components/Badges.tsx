import { docTypeNames, label, motherStatusNames, participantStatusNames, profileStatusNames, requestStatusNames } from '../utils/labels'

export function Badge({ text, kind }: { text: string; kind: 'ok' | 'wait' | 'muted' }) {
  return <span className={'badge badge-' + kind}>{text}</span>
}
export function ProfileBadge({ status }: { status: string }) {
  return <Badge text={label(profileStatusNames, status)} kind={status === 'checked' ? 'ok' : 'wait'} />
}
export function MotherStatusBadge({ status }: { status: string }) {
  const kind = status === 'active' ? 'ok' : status === 'inactive' ? 'muted' : 'wait'
  return <Badge text={label(motherStatusNames, status)} kind={kind} />
}
export function RequestBadge({ status }: { status: string }) {
  const kind = status === 'issued' ? 'ok' : status === 'rejected' ? 'muted' : 'wait'
  return <Badge text={label(requestStatusNames, status)} kind={kind} />
}
export function ParticipantBadge({ status }: { status: string }) {
  const kind = status === 'came' ? 'ok' : status === 'absent' ? 'muted' : 'wait'
  return <Badge text={label(participantStatusNames, status)} kind={kind} />
}
export function DocBadge({ status }: { status: string }) {
  return <Badge text={status === 'signed' ? 'Подписан' : 'Не подписан'} kind={status === 'signed' ? 'ok' : 'wait'} />
}
export function docTypeText(type: string) {
  return label(docTypeNames, type)
}
