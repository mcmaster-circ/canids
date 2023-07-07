const baseUrl = process.env.NEXT_PUBLIC_API_URL

interface DataViewParams {
  views: { uuid: string; name: string }[]
  params: {
    start: string
    end: string
    interval: number
    maxSize: number
    from: number
  }
}

export const getChartsData = async ({ params }: { params: DataViewParams }) => {
  const envUrls = params.views.map(
    ({ uuid }) =>
      baseUrl +
      '/data/?view=' +
      uuid +
      '&' +
      new URLSearchParams(params.params as any).toString()
  )
  const res = await Promise.all(
    envUrls.map((url) =>
      fetch(url, {
        method: 'GET',
        credentials: 'include',
      })
    )
  )
  const data = await Promise.all(res.map((r) => r.json()))
  return data.map((c, i) => ({
    ...c,
    uuid: params.views[i].uuid,
    name: params.views[i].name,
  }))
}
