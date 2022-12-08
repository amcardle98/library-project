import './App.css'
import useSWR, { KeyedMutator } from 'swr'
import { Box, List } from '@mantine/core'
import { useState } from 'react'
import { Input, Card, CardContent, Button, IconButton } from '@mui/material'
import { useForm } from '@mantine/form'
import { create } from '@mui/material/styles/createTransitions'
import AddBook from './components/AddBook'
import FindBook from './components/FindBook'
//import DeleteBook from './components/DeleteBook'
import { Delete } from '@mui/icons-material'

export interface Book {
  id: number
  title: string
  done: boolean
}

export const ENDPOINT = "http://localhost:4000"

// get the books and store them in an array

const fetcher = (url: string) => fetch(`${ENDPOINT}/${url}`).then((res) => res.json())

//create a empty array to store the books
const books: Book[] = []

function deleteBook(id: number, mutate: KeyedMutator<any>, book: Book[]) {
  return async () => {
    await fetch(`${ENDPOINT}/api/delete_book/${id}`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id })
    })

    mutate({ data: book.filter((book) => book.id !== id) }, false)
  }
}

function App() {
  const { data, mutate } = useSWR<{data: Book[]}>('api/books', fetcher)

  return <Card className='container' variant='outlined' sx={{maxWidth: 345}}>
    <div className='card'>
      <CardContent >
        <h3>Library</h3>  
        <div className='buttons'>
          <div className='add'> 
            <AddBook mutate={mutate} /> 
          </div>
          <div className='find'>
          <FindBook mutate={mutate} />
          </div>
        </div>
      </CardContent>
    </div>
    {/* <ViewBook mutate={mutate} /> */}
    <div className='library'>
      <List spacing="xs" size="sm" mb={12} center>
        {data?.data?.map((book) => (
          <Card key={book.id}>
            <CardContent className='card'>
              <h3>{book.title}</h3>
            </CardContent>
            <IconButton onClick={deleteBook(book.id, mutate, data?.data)}>
              <Delete />
            </IconButton>
          </Card>
        ))}
      </List>
    </div>
  </Card>
}

export default App
