export interface BaseModal {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: never | null
}

export interface EnumModal {
  value: number,
  label: string,
  color: string
}
