import { createTheme } from '@mantine/core';

export const theme = createTheme({
    colors: {
        // Kurumsal lacivert tonlarına göre oluşturulmuş renk paleti
        meuBlue: [
            '#eef1f6', '#d5dbe6', '#a5b4cc', '#738cb3', '#49689d',
            '#2d508f', '#1d4287', '#123473', '#0e2e68', '#07265c'
        ],
        // Kurumsal turuncu tonlarına göre oluşturulmuş renk paleti
        meuOrange: [
            '#fff0e6', '#ffe0cc', '#ffc199', '#ff9f66', '#ff813a',
            '#ff6a1e', '#ff5a09', '#e44800', '#cb3f00', '#b13600'
        ],
    },
    primaryColor: 'meuBlue', // Sistemin ana rengi lacivert
    primaryShade: 8, // Koyu ve kurumsal bir ton seçimi
    fontFamily: 'Inter, sans-serif',
    components: {
        Button: {
            defaultProps: {
                color: 'meuOrange', // Butonlar dikkat çekici kurumsal turuncu olacak
            },
        },
        Card: {
            styles: {
                root: {
                    backgroundColor: '#ffffff',
                    borderColor: '#d5dbe6', // Hafif mavi-gri kart sınırları
                }
            }
        }
    },
});