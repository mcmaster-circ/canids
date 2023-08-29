import Image from 'next/image'
import mailgun from '@images/mail_logos/mailgun.png'
import sendgrid from '@images/mail_logos/sendgrid.png'
import postal from '@images/mail_logos/postal.svg'
import postmark from '@images/mail_logos/postmark.png'
import sparkpost from '@images/mail_logos/sparkpost.png'
import { Cancel } from '@mui/icons-material'

export interface SectionProps {
  data: any
  setData: (a: any) => void
  initialData: any
}
export interface FormProps {
  service: string
  url: string
  apiKey: string
  fromAddress: string
  fromName: string
  domain: string
  accessURL: string
}

export const defaultFormValues = {
  service: '',
  url: '',
  apiKey: '',
  fromAddress: '',
  fromName: '',
  domain: '',
  accessURL: '',
}

export const steps = ['Service', 'Configuration']

export const typeButtons = [
  {
    label: 'None',
    icon: <Cancel />,
    key: 'NONE',
  },
  {
    label: 'Mailgun',
    icon: <Image src={mailgun} alt={'Mailgun'} priority={true} height={20} />,
    key: 'MAILGUN',
  },
  {
    label: 'SendGrid',
    icon: <Image src={sendgrid} alt={'SendGrid'} priority={true} height={20} />,
    key: 'SENDGRID',
  },
  {
    label: 'Postal',
    icon: <Image src={postal} alt={'Postal'} priority={true} height={20} />,
    key: 'POSTAL',
  },
  {
    label: 'Postmark',
    icon: <Image src={postmark} alt={'Postmark'} priority={true} height={20} />,
    key: 'POSTMARK',
  },
  {
    label: 'SparkPost',
    icon: (
      <Image src={sparkpost} alt={'SparkPost'} priority={true} height={20} />
    ),
    key: 'SPARKPOST',
  },
]
