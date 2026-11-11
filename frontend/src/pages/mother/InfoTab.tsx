import ProfileCard from './ProfileCard'
import StatusCard from './StatusCard'
import ChildrenCard from './ChildrenCard'
import SecretCard from './SecretCard'
import type { Child, MotherStatusInfo, Profile } from '../../types'

type Props = { motherId: number; profile: Profile; status: MotherStatusInfo; childList: Child[]; onChanged: () => void }

export default function InfoTab({ motherId, profile, status, childList, onChanged }: Props) {
  return (
    <div className="grid-2">
      <div>
        <ProfileCard motherId={motherId} profile={profile} onChanged={onChanged} />
        <ChildrenCard motherId={motherId} items={childList} onChanged={onChanged} />
      </div>
      <div>
        <StatusCard motherId={motherId} status={status} onChanged={onChanged} />
        <SecretCard motherId={motherId} />
      </div>
    </div>
  )
}
